package controllers

import (
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/kOstovic/CoffeeMachine/cmd/CoffeeMachineRestApiV3/repository"
)

type Ingredient struct {
	Water       uint16 `form:"Water" json:"Water"`
	Milk        uint16 `form:"Milk" json:"Milk"`
	Sugar       uint16 `form:"Sugar" json:"Sugar"`
	CoffeeBeans uint16 `form:"CoffeeBeans" json:"CoffeeBeans"`
	TeaBeans    uint16 `form:"TeaBeans" json:"TeaBeans"`
	Cups        uint16 `form:"Cups" json:"Cups"`
}

// register route for /ingredient in gin framework
func RegisterRoutesIngredients(router *gin.RouterGroup) {
	router.GET("", getIngredients)
	router.PUT("", putIngredients)
	router.PATCH("", patchIngredients)
}

func getIngredients(c *gin.Context) {
	if c.Query("name") != "" {
		getIngredientsByName(c)
	} else {
		getAllIngredients(c)
	}
}
func putIngredients(c *gin.Context) {
	if c.Query("name") != "" {
		putIngredientsByName(c)
	} else {
		putAllIngredients(c)
	}
}

// getAllIngredients godoc
// @Summary Get all ingredients available
// @Description Get all ingredients available
// @Produce json
// @Success 200 {array} Ingredient
// @Failure 400,401,404
// @Router /ingredients [get]
// @Security BearerAuth
func getAllIngredients(c *gin.Context) {
	if code, err := parseToken(c); err != nil {
		c.JSON(code, err.Error())
		return
	}
	c.JSON(http.StatusOK, repository.GetMachineIngredients())
}

// getIngredientsByName godoc
// @Summary Get ingredient by name from query
// @Description Get ingredient by name from query
// @Param name query string false "name of ingredient to get"
// @Produce json
// @Success 200 {object} Ingredient
// @Failure 400,401,404
// @Router /ingredients [get]
// @Security BearerAuth
func getIngredientsByName(c *gin.Context) {
	if code, err := parseToken(c); err != nil {
		c.JSON(code, err.Error())
		return
	}
	name := c.Query("name")
	cm, err := repository.GetIngredientValueByName(name)
	if err != nil {
		c.JSON(checkErrCode(err), err.Error())
		return
	} else {
		c.JSON(http.StatusOK, cm)
	}
}

// putAllIngredients godoc
// @Summary Update ingredients based on given Ingredient json, updates all
// @Description Update ingredients based on given Ingredient json, updates all
// @Param ingredient body Ingredient false "Update Ingredient object with Put option"
// @Accept json
// @Produce json
// @Success 200 {object} Ingredient
// @Failure 400,401,404
// @Failure 500
// @Router /ingredients [put]
// @Security BearerAuth
func putAllIngredients(c *gin.Context) {
	if code, err := parseToken(c); err != nil {
		c.JSON(code, err.Error())
		return
	}
	cm, err := checkIngredientsFromBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not update ingredients %v because of error: "+err.Error(), cm)
		return
	}
	cm, err = repository.UpdateIngredientPut(cm)
	if err != nil {
		c.JSON(checkErrCode(err), err.Error())
		log.Warnf("Could not update ingredients %v because of error: "+err.Error(), cm)
		return
	} else {
		log.Infof("Update Ingredients in machine to value %v", cm)
		c.JSON(http.StatusOK, cm)
	}

}

// putIngredientsByName godoc
// @Summary Update ingredients based on given Ingredient name and value in query or update all from body
// @Description Update ingredients based on given Ingredient name and value in query or update all from body
// @Param name query string false "name of ingredient to change"
// @Param value query int false "value of ingredient to change"
// @Param ingredient body Ingredient false "Update Ingredient object with Put option"
// @Produce json
// @Success 200 {object} Ingredient
// @Failure 400,401,404
// @Failure 500
// @Router /ingredients [put]
// @Security BearerAuth
func putIngredientsByName(c *gin.Context) {
	if code, err := parseToken(c); err != nil {
		c.JSON(code, err.Error())
		return
	}
	name := c.Query("name")
	valueStr := c.Query("value")
	if name == "" || valueStr == "" {
		c.JSON(http.StatusBadRequest, "name and value must be in query for putIngredientsByName operation")
		log.Warnf("Could not put ingredient by name because name is empty")
		return
	}
	value, err := strconv.ParseUint(valueStr, 10, 16)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not put in Ingredient called %v because of error: %v", name, err.Error())
		return
	}
	valueuint16 := uint16(value)
	cm, err := repository.PutIngredientsByName(name, valueuint16)
	if err != nil {
		c.JSON(checkErrCode(err), err.Error())
		log.Warnf("Could not put in Ingredient called %v because of error: %v", name, err.Error())
		return
	}
	log.Infof("Ingredient updated over put by name with values %v", cm)
	c.JSON(http.StatusOK, cm)
}

// patchIngredients godoc
// @Summary Update ingredients based on given Ingredient json, update only given
// @Description Update ingredients based on given Ingredient json, update only given
// @Param ingredient body Ingredient true "Update Ingredient object with Patch option"
// @Accept json
// @Produce json
// @Success 200 {object} Ingredient
// @Failure 400,401,404
// @Failure 500
// @Router /ingredients [patch]
// @Security BearerAuth
func patchIngredients(c *gin.Context) {
	if code, err := parseToken(c); err != nil {
		c.JSON(code, err.Error())
		return
	}
	cm, err := checkIngredientsFromBody(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not patch ingredient %v because of error: %v", err.Error(), cm)
		return
	}
	cm, err = repository.UpdateIngredientPatch(cm)
	if err != nil {
		c.JSON(checkErrCode(err), err.Error())
		log.Warnf("Could not patch ingredient %v because of error: %v", err.Error(), cm)
		return
	} else {
		c.JSON(http.StatusOK, cm)
		log.Infof("Ingredient updated over patch with values %v", cm)
		return
	}
}
