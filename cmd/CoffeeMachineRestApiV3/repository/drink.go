package repository

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/kOstovic/CoffeeMachine/internal/models"
)

func GetAvailableDrinks() ([]models.DrinkDB, error) {
	var allDrinks []models.DrinkDB
	if err := DbDrinks.Find(&allDrinks).Error; err != nil {
		log.Errorf("Failed to retrieve drinks: %s", err)
		return []models.DrinkDB{}, fmt.Errorf("Failed to retrieve drinks from database")
	}
	return allDrinks, nil
}

func GetDrinkByName(name string) (models.Drink, error) {
	var drinkDB models.DrinkDB
	if err := DbDrinks.Where("name = ? AND tenant_name = ?", name, "").First(&drinkDB).Error; err != nil {
		log.Errorf("Failed to retrieve drinks: %s", err)
		return models.Drink{}, fmt.Errorf("Failed to retrieve drinks from database")
	}
	return drinkDB.ConvertDrinkDBToDrink(), nil
}

func AddDrink(name string, drink models.Drink) (models.DrinkDB, error) {
	validation, err := drink.ValidationDrink()
	if !validation {
		return drink.ConvertDrinkToDB(name, ""), err
	}

	drinkDB := drink.ConvertDrinkToDB(name, "")
	dbresultDrink := DbDrinks.Create(&drinkDB)
	if dbresultDrink.Error != nil {
		log.Errorf("Could not save drink to drink database: " + dbresultDrink.Error.Error())
		return drink.ConvertDrinkToDB(name, ""), fmt.Errorf("Could not save drink to drink, drink already exists")
	}
	log.Debugf("Repository layer added drink with name: %s", name)
	return drinkDB, nil
}

func RemoveDrink(name string) (bool, error) {
	var drinkDB models.DrinkDB
	if err := DbDrinks.Where("name = ? AND tenant_name = ?", name, "").First(&drinkDB).Error; err != nil {
		log.Errorf("Removal - Failed to retrieve drink: %s", err)
		return false, fmt.Errorf("Drink does not exist therfor cannot be removed")
	}

	dbresultDrink := DbDrinks.Where("name = ? AND tenant_name = ?", name, "").Unscoped().Delete(&models.DrinkDB{})
	if dbresultDrink.Error != nil {
		log.Errorf("Could not delete drink from drink database: " + dbresultDrink.Error.Error())
		return false, fmt.Errorf("Could not delete drink from drink database: " + dbresultDrink.Error.Error())
	}
	log.Debugf("Repository layer removed drink with name: %s", name)
	return true, nil
}

func DeactivateDrink(name string) (bool, error) {
	var drinkDB models.DrinkDB
	if err := DbDrinks.Where("name = ? AND tenant_name = ?", name, "").First(&drinkDB).Error; err != nil {
		log.Errorf("Deactivation - Failed to retrieve drink: %s", err)
		return false, fmt.Errorf("Drink does not exist therfor cannot be deactivated")
	}

	dbresultDrink := DbDrinks.Where("name = ? AND tenant_name = ?", name, "").Delete(&models.DrinkDB{})
	if dbresultDrink.Error != nil {
		log.Errorf("Could not deactivate drink from drink database: " + dbresultDrink.Error.Error())
		return false, fmt.Errorf("Could not deactivate drink from drink database: " + dbresultDrink.Error.Error())
	}
	log.Debugf("Repository layer deactivated drink with name: %s", name)
	return true, nil
}

func ActivateDrink(name string) (bool, error) {
	var drinkDB models.DrinkDB
	if err := DbDrinks.Unscoped().Where("name = ? AND tenant_name = ?", name, "").First(&drinkDB).Error; err != nil {
		log.Errorf("Activation -  Failed to retrieve drink: %s", err)
		return false, fmt.Errorf("Drink does not exist therfor cannot be activated")
	}
	if err := DbDrinks.Where("name = ? AND tenant_name = ?", name, "").First(&drinkDB).Error; err == nil {
		log.Errorf("Activation - Drink %s already activate", name)
		return false, fmt.Errorf("Activation - Drink %s already activate", name)
	}

	dbresultDrink := DbDrinks.Unscoped().Where("name = ? AND tenant_name = ?", name, "").Model(&drinkDB).Update("deleted_at", nil)
	if dbresultDrink.Error != nil {
		log.Errorf("Could not activate drink from drink database: " + dbresultDrink.Error.Error())
		return false, fmt.Errorf("Could not activate drink from drink database")
	}
	log.Debugf("Repository layer deactivated drink with name: %s", name)
	return true, nil
}

func GetConsumeDrink(name string, den models.Denomination) (bool, models.Denomination, error) {
	var drinkDB models.DrinkDB
	var denRet models.Denomination
	if err := DbDrinks.Where("name = ? AND tenant_name = ?", name, "").First(&drinkDB).Error; err != nil {
		log.Errorf("Failed to retrieve drinks: %s", err)
		return false, den, fmt.Errorf("Drink does not exist")
	}

	resourcePrereq, cost, newIng, err := drinkDB.CheckResourcePrereqForDrinkDB()
	if err != nil {
		log.Debugf("Cannot not consume drink called %v because of error: %v: "+err.Error(), name)
		return false, den, err
	}
	if !resourcePrereq {
		log.Debugf("CoffeeMachine does not have enough resources for %v drink", name)
		return false, den, nil
	}

	newDen, err := models.CalculateDenominationAfterConsume(den, cost)
	if err != nil {
		log.Debugf("Cannot not consume drink called %v because of error: "+err.Error(), name)
		return false, den, err
	}

	txErr := DbDenomination.Transaction(func(txDenomination *gorm.DB) error {
		newDenDb := newDen.ConvertDenominationToDenominationDB("")
		dbresultDen := txDenomination.Where("tenant_name = ?", "").Model(&newDenDb).Updates(map[string]interface{}{"Half": newDenDb.Half, "One": newDenDb.One, "Two": newDenDb.Two, "Five": newDenDb.Five, "Ten": newDenDb.Ten, "Total": newDenDb.Total})
		if dbresultDen.Error != nil {
			log.Errorf("Could not save denomination to denomination database in consume action: " + dbresultDen.Error.Error())
			return dbresultDen.Error
		}
		err := DbIngredient.Transaction(func(txIngredient *gorm.DB) error {
			newIngDb := newIng.ConvertIngredientToIngredientDB("")
			dbresultIng := txIngredient.Where("tenant_name = ?", "").Model(&newIngDb).Updates(map[string]interface{}{"Water": newIngDb.Water, "Milk": newIngDb.Milk, "Sugar": newIngDb.Sugar, "CoffeeBeans": newIngDb.CoffeeBeans, "TeaBeans": newIngDb.TeaBeans, "Cups": newIngDb.Cups})
			if dbresultIng.Error != nil {
				log.Errorf("Could not save ingredient to ingredient database: " + dbresultIng.Error.Error())
				return dbresultIng.Error
			}
			denRet, err = models.UpdateDenominationConsume(den, cost)
			if err != nil {
				log.Debugf("Could not consume drink called %v because of error: %v", name, err.Error())
				return err
			}
			check, err := drinkDB.ConsumeDrinkDB()
			if !check || err != nil {
				log.Debugf("Could not consume drink called %v because of error: %v", name, err.Error())
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		log.Debugf("Could not consume drink called %v because of error: %v", name, txErr.Error())
		return false, den, fmt.Errorf("Could not consume drink called %v because of error: %v", name, txErr.Error())
	}

	log.Debugf("Consumed drink called %v saved in repository layer", name)
	return true, denRet, nil
}
