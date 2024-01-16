package repository

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/kOstovic/CoffeeMachine/internal/models"
)

func GetMachineIngredients() *models.Ingredient {
	return models.GetMachineIngredients()
}

func GetIngredientValueByName(name string) (string, error) {
	return models.GetIngredienteValueByName(name)
}

func UpdateIngredientPut(ing models.Ingredient) (models.Ingredient, error) {
	if !MachineInitialized {
		return models.Ingredient{}, fmt.Errorf("Machine must be initialized to use this endpoint")
	}
	validation, err := ing.ValidationIngredient()
	if !validation {
		log.Errorf("Error occurred: " + err.Error())
		return models.Ingredient{}, err
	}
	ingDb := ing.ConvertIngredientToIngredientDB("")
	dbresultIng := DbIngredient.Where("tenant_name = ?", "").Model(&ingDb).Updates(map[string]interface{}{"Water": ingDb.Water, "Milk": ingDb.Milk, "Sugar": ingDb.Sugar, "CoffeeBeans": ingDb.CoffeeBeans, "TeaBeans": ingDb.TeaBeans, "Cups": ingDb.Cups})
	if dbresultIng.Error != nil {
		log.Errorf("Could not save ingredient to ingredient database: " + dbresultIng.Error.Error())
		return models.Ingredient{}, fmt.Errorf("Could not save ingredient to ingredient database")
	}
	_, err = models.UpdateIngredientPut(ing)
	if err != nil {
		log.Errorf("Could not update ingredients because of error: " + err.Error())
		return ing, err
	}

	return ing, nil
}

func PutIngredientsByName(ingredient string, value uint16) (models.Ingredient, error) {
	if !MachineInitialized {
		return models.Ingredient{}, fmt.Errorf("Machine must be initialized to use this endpoint")
	}

	ingDb := models.GetMachineIngredients().ConvertIngredientToIngredientDB("")
	dbresultIng := DbIngredient.Where("tenant_name = ?", "").Model(&ingDb).Update(ingredient, value)
	if dbresultIng.Error != nil {
		log.Errorf("Could not save ingredient to ingredient database: " + dbresultIng.Error.Error())
		return models.Ingredient{}, fmt.Errorf("Could not save ingredient to ingredient database")
	}
	cm, err := models.UpdateIngredientValueByName(ingredient, value)
	if err != nil {
		log.Warnf("Could not put in ingredient called %v because of error: %v", ingredient, err.Error())
		return models.Ingredient{}, err
	}

	log.Debugf("Update Ingredients in machine named %v to value %v", ingredient, value)
	return cm, nil
}

func UpdateIngredientPatch(ing models.Ingredient) (models.Ingredient, error) {
	if !MachineInitialized {
		return models.Ingredient{}, fmt.Errorf("Machine must be initialized to use this endpoint")
	}
	validation, err := ing.ValidationIngredient()
	if !validation {
		log.Debug("Repository layer validation of Ingredient failed")
		return models.Ingredient{}, err
	}

	// this would update all fields including zero value ones
	/*patchobj := make(map[string]interface{})
	if ing.Water > 0 {
		patchobj["Water"] = ing.Water
	}
	if ing.Milk > 0 {
		patchobj["Milk"] = ing.Milk
	}
	if ing.Sugar > 0 {
		patchobj["Sugar"] = ing.Sugar
	}
	if ing.CoffeeBeans > 0 {
		patchobj["CoffeeBeans"] = ing.CoffeeBeans
	}
	if ing.TeaBeans > 0 {
		patchobj["TeaBeans"] = ing.TeaBeans
	}
	if ing.Cups > 0 {
		patchobj["Cups"] = ing.Cups
	}
	dbresultIng := DbIngredient.Model(models.GetMachineIngredients().ConvertIngredientToIngredientDB(0,"")).Select(*).Updates(&patchobj)*/

	// this will update only non-zero fields
	ingDb := ing.ConvertIngredientToIngredientDB("")
	modelDb := models.GetMachineIngredients().ConvertIngredientToIngredientDB("")
	dbresultIng := DbIngredient.Where("tenant_name = ?", "").Model(&modelDb).Updates(&ingDb)
	if dbresultIng.Error != nil {
		log.Errorf("Could not save ingredient to ingredient database: " + dbresultIng.Error.Error())
		return models.Ingredient{}, fmt.Errorf("Could not save ingredient to ingredient database")
	}

	cm, err := models.UpdateIngredientPatch(ing)
	if err != nil {
		log.Errorf("Could not patch ingredient because of error: %v", err.Error())
		return cm, err
	}

	log.Debugf("Patched Ingredients in machine %v", cm)
	return cm, nil
}
