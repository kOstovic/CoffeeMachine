package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/kOstovic/CoffeeMachine/internal/models"
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
// @Success 200 {object} models.Denomination
// @Router /money [get]
func getAllAvailableDenomination(c *gin.Context) {
	c.JSON(http.StatusOK, models.GetCurrentMoney())
}

// getDenominationByName godoc
// @Summary Get denominations by name from query
// @Description Get denominations by name from query
// @Param name query string false "name of denomination to get"
// @Produce json
// @Success 200 {object} models.Denomination
// @Failure 400,404
// @Router /money [get]
func getDenominationByName(c *gin.Context) {
	name := c.Query("name")
	cm, err := models.GetDenominationValueByName(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
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
// @Success 200 {object} models.Denomination
// @Failure 400,404
// @Failure 500
// @Router /money [put]
func putAllDenomination(c *gin.Context) {
	cm, err := checkMoneyFromURL(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not update Denomination because of error: " + err.Error())
		return
	}
	cm, err = models.UpdateDenominationPut(cm)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not update Denomination because of error: " + err.Error())
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
// @Success 200 {object} models.Denomination
// @Failure 400,404
// @Failure 500
// @Router /money [put]
func putDenominationByName(c *gin.Context) {
	name := c.Query("name")
	valueStr := c.Query("value")
	if name == "" || valueStr == "" {
		c.JSON(http.StatusBadRequest, "name and value must be in query for putDenominationByName operation")
		log.Warnf("Could not put drink Denomination by name because name is empty")
		return
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not put in Denomination called %v because of error: %v", name, err.Error())
		return
	}
	cm, err := models.UpdateDenominationValueByName(name, value)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not put in Denomination called %v because of error: %v", name, err.Error())
		return
	}
	c.JSON(http.StatusOK, cm)
	log.Debugf("Update Denomination in machine named %v to value %v", name, value)
}

// patchDenomination godoc
// @Summary Update ingredients based on given Denomination json, update only given
// @Description Update ingredients based on given Denomination json, update only given
// @Param denomination body Denomination true "Update Denomination object with Patch option"
// @Accept json
// @Produce json
// @Success 200 {object} models.Ingredient
// @Failure 400,404
// @Failure 500
// @Router /money [patch]
func patchDenomination(c *gin.Context) {
	cm, err := checkMoneyFromURL(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not patch Denomination because of error: %v", err.Error())
		return
	}
	cm, err = models.UpdateDenominationPatch(cm)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not patch Denomination because of error: %v", err.Error())
		return
	} else {
		c.JSON(http.StatusOK, cm)
	}
	log.Debugf("Patched Denomination in machine %v", cm)
}

func checkMoneyFromURL(c *gin.Context) (models.Denomination, error) {
	var denomination Denomination
	if c.ShouldBindQuery(&denomination) == nil {
		log.Debugf("====== Bind By Query String ====== from request %v", denomination)
		return models.Denomination{Half: denomination.Half,
			One: denomination.One, Two: denomination.Two,
			Five: denomination.Five, Ten: denomination.Ten}, nil
	} else if c.ShouldBindJSON(&denomination) == nil {
		log.Debugf("====== Bind By JSON ====== from request %v", denomination)
		return models.Denomination{Half: denomination.Half,
			One: denomination.One, Two: denomination.Two,
			Five: denomination.Five, Ten: denomination.Ten}, nil
	} else {
		log.Debugf("Denomination could not be parsed from request %v", denomination)
		return models.Denomination{}, fmt.Errorf("Denomination could not be parsed")
	}
}
