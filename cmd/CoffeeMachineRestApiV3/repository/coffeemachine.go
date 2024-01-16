package repository

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/kOstovic/CoffeeMachine/internal/models"
)

var (
	MachineInitialized bool = false
)

func InitializeMachine(ing models.Ingredient, den models.Denomination) (bool, error) {
	if MachineInitialized == true {
		log.Errorf("coffeeMachine cannot be initialized more than once")
		return false, fmt.Errorf("coffeeMachine cannot be initialized more than once")
	}

	validationIng, err := ing.ValidationIngredient()
	if !validationIng {
		log.Errorf("Error occurred: " + err.Error())
		return false, err
	}

	validationDen, err := den.ValidationDenomination()
	if !validationDen {
		log.Errorf("Error occurred: " + err.Error())
		return false, err
	}

	txErr := DbIngredient.Transaction(func(txIngredient *gorm.DB) error {
		ingDB := ing.ConvertIngredientToIngredientDB("")
		if err := txIngredient.Create(&ingDB).Error; err != nil {
			log.Errorf("Could not save ingredient to ingredient database: " + err.Error())
			return err
		}

		denDB := den.ConvertDenominationToDenominationDB("")
		err := DbDenomination.Transaction(func(txDenomination *gorm.DB) error {
			if err := txDenomination.Create(&denDB).Error; err != nil {
				log.Errorf("Could not save denomination to denomination database: " + err.Error())
				return err
			}

			_, errIng := models.InitializeIngredients(ing)
			if errIng != nil {
				log.Errorf("coffeeMachine could not be initialized as there was error in inmemory save for ingredient: " + errIng.Error())
				return errIng
			}
			_, errDen := models.InitializeDenominations(den)
			if errDen != nil {
				log.Errorf("coffeeMachine could not be initialized as there was error in inmemory save for denomination: " + errDen.Error())
				return errDen
			}

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		log.Errorf("coffeeMachine could not be initialized as there was error: " + txErr.Error())
		return false, fmt.Errorf("coffeeMachine could not be initialized as there was error: " + txErr.Error())
	}
	MachineInitialized = true
	log.Debugf("coffeeMachine initialized in repository layer with following parameters: Ingredients: %v Money: %v", ing, den)
	return true, nil
}

func DeleteDeInitializeMachine() (bool, error) {
	if MachineInitialized == false {
		log.Errorf("coffeeMachine cannot be deinitialized more than once")
		return false, fmt.Errorf("coffeeMachine cannot be deinitialized more than once")
	}

	txErr := DbIngredient.Transaction(func(txIngredient *gorm.DB) error {
		ingDb := models.GetMachineIngredients().ConvertIngredientToIngredientDB("")
		dbresultIng := txIngredient.Where("tenant_name = ?", "").Model(&ingDb).Unscoped().Delete(&ingDb)
		if dbresultIng.Error != nil {
			log.Errorf("Could not save ingredient to ingredient database: " + dbresultIng.Error.Error())
			return dbresultIng.Error
		}

		err := DbDenomination.Transaction(func(txDenomination *gorm.DB) error {
			denDb := models.GetCurrentMoney().ConvertDenominationToDenominationDB("")
			dbresultDen := txDenomination.Where("tenant_name = ?", "").Model(&denDb).Unscoped().Delete(&denDb)
			if dbresultDen.Error != nil {
				log.Errorf("Could not save denomination to denomination database: " + dbresultDen.Error.Error())
				return dbresultDen.Error
			}

			_, errIng := models.CleanupIngredients()
			if errIng != nil {
				log.Errorf("coffeeMachine could not be DeInitialized because of error: " + errIng.Error())
				return errIng
			}
			_, errDen := models.CleanupDenominations()
			if errDen != nil {
				log.Errorf("coffeeMachine could not be DeInitialized because of error: " + errIng.Error())
				return errDen
			}

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		log.Errorf("coffeeMachine could not be Deinitialized as there was error: " + txErr.Error())
		return false, fmt.Errorf("coffeeMachine could not be Deinitialized as there was error: " + txErr.Error())
	}
	MachineInitialized = false
	log.Debugf("coffeeMachine Deinitialized in repository layer")
	return true, nil
}
