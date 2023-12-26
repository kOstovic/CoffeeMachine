package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/kOstovic/CoffeeMachine/internal/models"
)

type Drink struct {
	Water       uint16  `form:"Water" json:"Water" binding:"required_without_all=Milk Sugar CoffeeBeans TeaBeans Cups Money"`
	Milk        uint16  `form:"Milk" json:"Milk" binding:"required_without_all=Water Sugar CoffeeBeans TeaBeans Cups Money"`
	Sugar       uint16  `form:"Sugar" json:"Sugar" binding:"required_without_all=Milk Water CoffeeBeans TeaBeans Cups Money"`
	CoffeeBeans uint16  `form:"CoffeeBeans" json:"CoffeeBeans" binding:"required_without_all=Milk Sugar Water TeaBeans Cups Money"`
	TeaBeans    uint16  `form:"TeaBeans" json:"TeaBeans" binding:"required_without_all=Milk Sugar CoffeeBeans Water Cups Money"`
	Cups        uint16  `form:"Cups" json:"Cups" binding:"required_without_all=Milk Sugar CoffeeBeans TeaBeans Water Money"`
	Money       float64 `form:"Money" json:"Money" binding:"required_without_all=Milk Sugar CoffeeBeans TeaBeans Cups Water"`
}

// register route for drink in gin framework
func RegisterRoutesDrink(router *gin.RouterGroup) {
	router.GET("", getDrink)
	router.GET("/consume", getConsumeDrink)
	router.POST("", postAddDrink)
	router.DELETE("", postRemoveDrink)
}
func getDrink(c *gin.Context) {
	if c.Query("name") != "" {
		getConsumeDrink(c)
	} else {
		getAllAvailableDrinks(c)
	}
}

// getAllAvailableDrinks godoc
// @Summary Get all drinks available
// @Description Get all drinks available
// @Produce application/json
// @Success 200 {array} models.Drink
// @Router /drinks [get]
func getAllAvailableDrinks(c *gin.Context) {
	c.JSON(http.StatusOK, models.GetAvailableDrinks())
}

// ConsumeDrink godoc
// @Summary Consumes Drink over
// @Description Consumes Drink over
// @Param name query string true "Name of drink to consume"
// @Param Half query string false "Denomination Half to consume"
// @Param One query string false "Denomination One to consume"
// @Param Two query string false "Denomination Two to consume"
// @Param Five query string false "Denomination Five to consume"
// @Param Ten query string false "Denomination Ten to consume"
// @Produce application/json
// @Success 200 {object} models.Denomination
// @Failure 400,404
// @Failure 500
// @Router /drinks/consume [get]
func getConsumeDrink(c *gin.Context) {
	name := c.Query("name")
	preReq, cost, err := models.CheckPrereqForDrink(name)

	if !preReq || err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Debugf("Could not consume drink called %v because of error: %v", name, err.Error())
		return
	} else {
		denominationParam, err := checkDenFromReq(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			log.Debugf("Could not consume drink called %v because of error: %v", name, err.Error())
			return
		}
		denRet, err := models.UpdateDenominationConsume(denominationParam, cost)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			log.Debugf("Could not consume drink called %v because of error: %v", name, err.Error())
			return
		} else {
			models.ConsumeDrink(name)
			c.JSON(http.StatusOK, denRet)
		}
	}
	log.Debugf("Consumed drink called %v", name)
}

// AddDrink godoc
// @Summary Initialize new drink to consume on given Drink json
// @Description Initialize new drink to consume on given Drink json
// @Param name query string true "name of drink to create"
// @Param Drink body Drink true "Add Drink object"
// @Accept  json
// @Produce json
// @Success 200 {object} models.Drink
// @Failure 400,404
// @Failure 500
// @Router /drinks [post]
func postAddDrink(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, "name in query is empty")
		log.Warnf("Could not add drink because name is empty")
		return
	}
	cm, err := checkDrinkFromURL(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not add drink called %v because of error: %v", name, err.Error())
		return
	}
	drink, err := models.AddDrink(name, cm)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not add drink called %v because of error: %v", name, err.Error())
		return
	} else {
		c.JSON(http.StatusOK, drink)
	}
	log.Infof("Added drink called %v", name)
}

// RemoveDrink godoc
// @Summary Remove drink from machine on given name
// @Description Remove drink from machine on given name
// @Param name query string true "name of drink to delete"
// @Produce json
// @Success 200 {object} bool
// @Failure 400,404
// @Failure 500
// @Router /drinks [delete]
func postRemoveDrink(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, "name in query is empty")
		log.Warnf("Could not remove drink because name is empty")
		return
	}
	OK, err := models.RemoveDrink(name)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Warnf("Could not remove drink called %v because of error: %v", name, err.Error())
		return
	} else {
		c.JSON(http.StatusOK, OK)
	}
	log.Infof("Removed drink called %v", name)
}

func checkDrinkFromURL(c *gin.Context) (models.Drink, error) {
	var drink Drink
	if c.ShouldBindQuery(&drink) == nil {
		log.Debugf("====== Bind By Query String ====== from request %v", &drink)
		return models.Drink{Water: drink.Water,
			Milk: drink.Milk, Sugar: drink.Sugar,
			CoffeeBeans: drink.CoffeeBeans, TeaBeans: drink.TeaBeans,
			Cups: drink.Cups, Money: drink.Money}, nil
	} else if c.ShouldBindJSON(&drink) == nil {
		log.Debugf("====== Bind By JSON ====== from request %v", &drink)
		return models.Drink{Water: drink.Water,
			Milk: drink.Milk, Sugar: drink.Sugar,
			CoffeeBeans: drink.CoffeeBeans, TeaBeans: drink.TeaBeans,
			Cups: drink.Cups, Money: drink.Money}, nil
	} else {
		log.Warnf("Drink could not be parsed from request %v", &drink)
		return models.Drink{}, fmt.Errorf("Drink could not be parsed in both query and body")
	}
}

func checkDenFromReq(c *gin.Context) (models.Denomination, error) {
	var money Money
	if c.ShouldBindQuery(&money) == nil {
		log.Debugf("====== Bind By Query String ====== from request %v", money)
		return models.Denomination{Half: money.Half,
			One: money.One, Two: money.Two,
			Five: money.Five, Ten: money.Ten}, nil
	} else if c.ShouldBindJSON(&money) == nil {
		log.Debugf("====== Bind By JSON ====== from request %v", money)
		return models.Denomination{Half: money.Half,
			One: money.One, Two: money.Two,
			Five: money.Five, Ten: money.Ten}, nil
	} else {
		log.Warnf("Denomination could not be parsed from request %v", money)
		return models.Denomination{}, fmt.Errorf("Denomination could not be parsed")
	}
}
