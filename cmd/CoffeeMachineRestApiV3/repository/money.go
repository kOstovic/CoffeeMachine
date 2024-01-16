package repository

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/kOstovic/CoffeeMachine/internal/models"
)

func GetCurrentMoney() *models.Denomination {
	return models.GetCurrentMoney()
}

func GetDenominationValueByName(denomination string) (string, error) {
	return models.GetDenominationValueByName(denomination)
}

func UpdateDenominationPut(den models.Denomination) (models.Denomination, error) {
	if !MachineInitialized {
		return models.Denomination{}, fmt.Errorf("Machine must be initialized to use this endpoint")
	}
	validation, err := den.ValidationDenomination()
	if !validation {
		log.Errorf("Error occurred: " + err.Error())
		return models.Denomination{}, err
	}

	denDb := den.ConvertDenominationToDenominationDB("")
	dbresultDen := DbDenomination.Where("tenant_name = ?", "").Model(&denDb).Updates(map[string]interface{}{"Half": denDb.Half, "One": denDb.One, "Two": denDb.Two, "Five": denDb.Five, "Ten": denDb.Ten, "Total": denDb.Total})
	if dbresultDen.Error != nil {
		log.Errorf("Could not save denomination to denomination database: " + dbresultDen.Error.Error())
		return models.Denomination{}, fmt.Errorf("Could not save ingredient to ingredient database")
	}
	_, err = models.UpdateDenominationPut(den)
	if err != nil {
		log.Errorf("Could not update denominations because of error: " + err.Error())
		return den, err
	}
	log.Debugf("Update Denomination in machine to value %v", den)
	return den, nil
}

func UpdateDenominationValueByName(denomination string, value int) (models.Denomination, error) {
	if !MachineInitialized {
		return models.Denomination{}, fmt.Errorf("Machine must be initialized to use this endpoint")
	}
	if value < 0 {
		return models.Denomination{}, fmt.Errorf("Value in Denomination cannot be negative'%v'", value)
	}

	denDb := models.GetCurrentMoney().ConvertDenominationToDenominationDB("")
	dbresultDen := DbDenomination.Where("tenant_name = ?", "").Model(&denDb).Update(denomination, value)
	if dbresultDen.Error != nil {
		log.Errorf("Could not save denomination to denomination database: " + dbresultDen.Error.Error())
		return models.Denomination{}, fmt.Errorf("Could not save ingredient to ingredient database")
	}
	dm, err := models.UpdateDenominationValueByName(denomination, value)
	if err != nil {
		log.Warnf("Could not put in denomination called %v because of error: %v", denomination, err.Error())
		return models.Denomination{}, err
	}

	log.Debugf("Update Denomination in machine named %v to value %v", denomination, value)
	return dm, nil
}

func UpdateDenominationPatch(den models.Denomination) (models.Denomination, error) {
	if !MachineInitialized {
		return models.Denomination{}, fmt.Errorf("Machine must be initialized to use this endpoint")
	}
	validation, err := den.ValidationDenomination()
	if !validation {
		log.Debug("Repository layer validation of Denomination failed")
		return models.Denomination{}, err
	}

	// this will update only non-zero fields
	denDb := den.ConvertDenominationToDenominationDB("")
	modelDb := models.GetCurrentMoney().ConvertDenominationToDenominationDB("")
	dbresultDen := DbDenomination.Where("tenant_name = ?", "").Model(&modelDb).Updates(&denDb)
	if dbresultDen.Error != nil {
		log.Errorf("Could not save denomination to denomination database: " + dbresultDen.Error.Error())
		return models.Denomination{}, fmt.Errorf("Could not save ingredient to ingredient database")
	}
	cm, err := models.UpdateDenominationPatch(den)
	if err != nil {
		log.Warnf("Could not patch Denomination because of error: %v", err.Error())
		return cm, err
	}
	log.Debugf("Patched Denomination in machine %v", cm)
	return cm, nil
}
