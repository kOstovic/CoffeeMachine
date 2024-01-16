package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/kOstovic/CoffeeMachine/cmd/CoffeeMachineRestApiV3/repository"
)

type Money struct {
	Half int `form:"Half" json:"Half"`
	One  int `form:"One" json:"One"`
	Two  int `form:"Two" json:"Two"`
	Five int `form:"Five" json:"Five"`
	Ten  int `form:"Ten" json:"Ten"`
}

// register route for coffeemachine in gin framework
func RegisterRoutesDenomination(router *gin.RouterGroup) {
	router.GET("", getDenominations)
	router.PUT("", putDenomination)
	router.PATCH("", patchDenomination)
}

func getDenominations(c *gin.Context) {
	if c.Query("name") != "" {
		getDenominationByName(c)
	} else {
		getAllAvailableDenomination(c)
	}
}
func putDenomination(c *gin.Context) {
	if c.Query("name") != "" {
		putDenominationByName(c)
	} else {
		putAllDenomination(c)
	}
}

// getAllAvailableDenomination godoc
// @Summary Get all denominations available
// @Description Get all denominations available
// @Produce json
// @Success 200 {object} Denomination
// @Router /money [get]
func getAllAvailableDenomination(c *gin.Context) {
	c.JSON(http.StatusOK, repository.GetCurrentMoney())
}

// getDenominationByName godoc
// @Summary Get denominations by name from query
// @Description Get denominations by name from query
// @Param name query string false "name of denomination to get"
// @Produce json
// @Success 200 {object} Denomination
// @Failure 400,404
// @Router /money [get]
func getDenominationByName(c *gin.Context) {
	name := c.Query("name")
	cm, err := repository.GetDenominationValueByName(name)
	if err != nil {
		c.JSON(checkErrCode(err), err.Error())
		return
	} else {
		c.JSON(http.StatusOK, cm)
	}
}

// putDenomination godoc
// @Summary Update Denomination based on given Denomination json, updates all
// @Description Update Denomination based on given Denomination json, updates all
// @Param denomination body Denomination true "Update Denomination object with Put option"
// @Accept json
// @Produce json
// @Success 200 {object} Denomination
// @Failure 400,404
// @Failure 500
// @Router /money [put]
func putAllDenomination(c *gin.Context) {
	cm, err := checkDenFromBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not update Denomination %v because of error: "+err.Error(), cm)
		return
	}
	cm, err = repository.UpdateDenominationPut(cm)
	if err != nil {
		c.JSON(checkErrCode(err), err.Error())
		log.Warnf("Could not update Denomination %v because of error: "+err.Error(), cm)
		return
	} else {
		c.JSON(http.StatusOK, cm)
	}
	log.Infof("Update Denomination in machine to value %v", cm)
}

// putDenominationByName godoc
// @Summary Update denomination based on given Denomination name and value in query or update all from body
// @Description Update denomination based on given Denomination name and value in query or update all from body
// @Param name query string false "name of denomination to change"
// @Param value query int false "value of denomination to change"
// @Param denomination body Denomination false "Update Denomination object with Put option"
// @Produce json
// @Success 200 {object} Denomination
// @Failure 400,404
// @Failure 500
// @Router /money [put]
func putDenominationByName(c *gin.Context) {
	name := c.Query("name")
	valueStr := c.Query("value")
	if name == "" || valueStr == "" {
		c.JSON(http.StatusBadRequest, "name and value must be in query for putDenominationByName operation")
		log.Warnf("Could not put Denomination by name because name is empty")
		return
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not put in Denomination called %v because of error: %v", name, err.Error())
		return
	}
	cm, err := repository.UpdateDenominationValueByName(name, value)
	if err != nil {
		c.JSON(checkErrCode(err), err.Error())
		log.Warnf("Could not put in Denomination called %v because of error: %v", name, err.Error())
		return
	}
	log.Infof("Denomination updated over put by name with values %v", cm)
	c.JSON(http.StatusOK, cm)

}

// patchDenomination godoc
// @Summary Update ingredients based on given Denomination json, update only given
// @Description Update ingredients based on given Denomination json, update only given
// @Param denomination body Denomination true "Update Denomination object with Patch option"
// @Accept json
// @Produce json
// @Success 200 {object} Ingredient
// @Failure 400,404
// @Failure 500
// @Router /money [patch]
func patchDenomination(c *gin.Context) {
	cm, err := checkDenFromBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not patch Denomination %v because of error: %v", err.Error(), cm)
		return
	}
	cm, err = repository.UpdateDenominationPatch(cm)
	if err != nil {
		c.JSON(checkErrCode(err), err.Error())
		log.Warnf("Could not patch Denomination %v because of error: %v", err.Error(), cm)
		return
	} else {
		c.JSON(http.StatusOK, cm)
		log.Infof("Denomination updated over patch with values %v", cm)
		return
	}
}
