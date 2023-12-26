package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/kOstovic/CoffeeMachine/internal/models"
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
// @Success 200 {array} models.Ingredient
// @Router /ingredients [get]
func getAllIngredients(c *gin.Context) {
	c.JSON(http.StatusOK, models.GetMachineIngredients())
}

// getIngredientsByName godoc
// @Summary Get ingredient by name from query
// @Description Get ingredient by name from query
// @Param name query string false "name of ingredient to get"
// @Produce json
// @Success 200 {object} models.Ingredient
// @Failure 400,404
// @Router /ingredients [get]
func getIngredientsByName(c *gin.Context) {
	name := c.Query("name")
	cm, err := models.GetIngredienteValueByName(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	} else {
		c.JSON(http.StatusOK, cm)
	}
}

// putIngredients godoc
// @Summary Update ingredients based on given Ingredient json, updates all
// @Description Update ingredients based on given Ingredient json, updates all
// @Param ingredient body Ingredient false "Update Ingredient object with Put option"
// @Accept json
// @Produce json
// @Success 200 {object} models.Ingredient
// @Failure 400,404
// @Failure 500
// @Router /ingredients [put]
func putAllIngredients(c *gin.Context) {
	cm, err := checkIngredientsFromURL(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not update ingredients because of error: " + err.Error())
		return
	}
	cm, err = models.UpdateIngredientPut(cm)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not update ingredients because of error: " + err.Error())
		return
	} else {
		c.JSON(http.StatusOK, cm)
	}
	log.Infof("Update Ingredients in machine to value %v", cm)
}

// putIngredientsByName godoc
// @Summary Update ingredients based on given Ingredient name and value in query or update all from body
// @Description Update ingredients based on given Ingredient name and value in query or update all from body
// @Param name query string false "name of ingredient to change"
// @Param value query int false "value of ingredient to change"
// @Param ingredient body Ingredient false "Update Ingredient object with Put option"
// @Produce json
// @Success 200 {object} models.Ingredient
// @Failure 400,404
// @Failure 500
// @Router /ingredients [put]
func putIngredientsByName(c *gin.Context) {
	name := c.Query("name")
	valueStr := c.Query("value")
	if name == "" || valueStr == "" {
		c.JSON(http.StatusBadRequest, "name and value must be in query for putIngredientsByName operation")
		log.Warnf("Could not put drink ingredient by name because name is empty")
		return
	}
	value, err := strconv.ParseUint(valueStr, 10, 16)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not put in ingredient called %v because of error: %v", name, err.Error())
		return
	}
	valueuint16 := uint16(value)
	cm, err := models.UpdateIngredientValueByName(name, valueuint16)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not put in ingredient called %v because of error: %v", name, err.Error())
		return
	}
	c.JSON(http.StatusOK, cm)
	log.Debugf("Update Ingredients in machine named %v to value %v", name, valueuint16)
}

// patchIngredients godoc
// @Summary Update ingredients based on given Ingredient json, update only given
// @Description Update ingredients based on given Ingredient json, update only given
// @Param ingredient body Ingredient true "Update Ingredient object with Patch option"
// @Accept json
// @Produce json
// @Success 200 {object} models.Ingredient
// @Failure 400,404
// @Failure 500
// @Router /ingredients [patch]
func patchIngredients(c *gin.Context) {
	cm, err := checkIngredientsFromURL(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not patch ingredient because of error: %v", err.Error())
		return
	}
	cm, err = models.UpdateIngredientPatch(cm)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not patch ingredient because of error: %v", err.Error())
		return
	} else {
		c.JSON(http.StatusOK, cm)
	}
	log.Debugf("Patched Ingredients in machine %v", cm)
}

func checkIngredientsFromURL(c *gin.Context) (models.Ingredient, error) {
	var ingredient Ingredient
	if c.ShouldBindQuery(&ingredient) == nil {
		log.Debugf("====== Bind By Query String ====== from request %v", ingredient)
		return models.Ingredient{Water: ingredient.Water,
			Milk: ingredient.Milk, Sugar: ingredient.Sugar,
			CoffeeBeans: ingredient.CoffeeBeans, TeaBeans: ingredient.TeaBeans,
			Cups: ingredient.Cups}, nil
	} else if c.ShouldBindJSON(&ingredient) == nil {
		log.Debugf("====== Bind By JSON ====== from request %v", ingredient)
		return models.Ingredient{Water: ingredient.Water,
			Milk: ingredient.Milk, Sugar: ingredient.Sugar,
			CoffeeBeans: ingredient.CoffeeBeans, TeaBeans: ingredient.TeaBeans,
			Cups: ingredient.Cups}, nil
	} else {
		log.Warnf("Ingredient could not be parsed from request %v", ingredient)
		return models.Ingredient{}, fmt.Errorf("Ingredient could not be parsed")
	}
}
