package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/kOstovic/CoffeeMachine/cmd/CoffeeMachineRestApiV3/repository"
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
	router.DELETE("", deleteRemoveDrink)
	router.POST("activate", postActivateDrink)
	router.DELETE("deactivate", deleteDeactivateDrink)
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
// @Success 200 {array} Drink
// @Router /drinks [get]
func getAllAvailableDrinks(c *gin.Context) {
	drinks, err := repository.GetAvailableDrinks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		log.Debugf("Could not get drinks because of error: %v", err.Error())
		return
	}
	c.JSON(http.StatusOK, drinks)
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
// @Success 200 {object} Denomination
// @Failure 400,404
// @Failure 500
// @Router /drinks/consume [get]
func getConsumeDrink(c *gin.Context) {
	name := c.Query("name")
	denominationParam, err := checkDenFromURL(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Debugf("Could not consume drink called %v because of error: %v", name, err.Error())
		return
	}

	check, denRet, err := repository.GetConsumeDrink(name, denominationParam)
	if !check || err != nil {
		c.JSON(checkErrCode(err), err.Error())
		log.Debugf("Could not consume drink called %v because of error: %v", name, err.Error())
		return
	} else {
		log.Infof("Consumed drink called %v", name)
		c.JSON(http.StatusOK, denRet)
	}
}

// AddDrink godoc
// @Summary Initialize new drink to consume on given Drink json
// @Description Initialize new drink to consume on given Drink json
// @Param name query string true "name of drink to create"
// @Param Drink body Drink true "Add Drink object"
// @Accept  json
// @Produce json
// @Success 200 {object} Drink
// @Failure 400,401,404
// @Failure 500
// @Router /drinks [post]
// @Security BearerAuth
func postAddDrink(c *gin.Context) {
	if code, err := parseToken(c); err != nil {
		c.JSON(code, err.Error())
		return
	}
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
	drinkDB, err := repository.AddDrink(name, cm)
	if err != nil {
		c.JSON(checkErrCode(err), err.Error())
		log.Warnf("Could not add drink called %v because of error: %v", name, err.Error())
		return
	} else {
		log.Infof("Added drink with params:  %v", drinkDB)
		c.JSON(http.StatusOK, drinkDB)
	}
}

// ActivateDrink godoc
// @Summary Activate drink from machine on given name
// @Description Activate drink from machine on given name
// @Param name query string true "name of drink to activate"
// @Produce json
// @Success 200 {object} bool
// @Failure 400,401,404
// @Failure 500
// @Router /drinks/activate [post]
// @Security BearerAuth
func postActivateDrink(c *gin.Context) {
	if code, err := parseToken(c); err != nil {
		c.JSON(code, err.Error())
		return
	}
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, "name in query is empty")
		log.Warnf("Could not activate drink because name is empty")
		return
	}
	OK, err := repository.ActivateDrink(name)
	if err != nil {
		c.JSON(checkErrCode(err), err.Error())
		log.Warnf("Could not activate drink called %v because of error: %v", name, err.Error())
		return
	} else {
		log.Infof("Activated drink called %v", name)
		c.JSON(http.StatusOK, OK)
	}
}

// DeactivateDrink godoc
// @Summary Deactivate drink from machine on given name
// @Description Deactivate drink from machine on given name
// @Param name query string true "name of drink to deactivate"
// @Produce json
// @Success 200 {object} bool
// @Failure 400,401,404
// @Failure 500
// @Router /drinks/deactivate [delete]
// @Security BearerAuth
func deleteDeactivateDrink(c *gin.Context) {
	if code, err := parseToken(c); err != nil {
		c.JSON(code, err.Error())
		return
	}
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, "name in query is empty")
		log.Warnf("Could not deactivate drink because name is empty")
		return
	}
	OK, err := repository.DeactivateDrink(name)
	if err != nil {
		c.JSON(checkErrCode(err), err.Error())
		log.Warnf("Could not deactivate drink called %v because of error: %v", name, err.Error())
		return
	} else {
		log.Infof("Deactivated drink called %v", name)
		c.JSON(http.StatusOK, OK)
	}
}

// RemoveDrink godoc
// @Summary Remove drink from machine on given name
// @Description Remove drink from machine on given name
// @Param name query string true "name of drink to delete"
// @Produce json
// @Success 200 {object} bool
// @Failure 400,401,404
// @Failure 500
// @Router /drinks [delete]
// @Security BearerAuth
func deleteRemoveDrink(c *gin.Context) {
	if code, err := parseToken(c); err != nil {
		c.JSON(code, err.Error())
		return
	}
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, "name in query is empty")
		log.Warnf("Could not remove drink because name is empty")
		return
	}
	OK, err := repository.RemoveDrink(name)
	if err != nil {
		c.JSON(checkErrCode(err), err.Error())
		log.Warnf("Could not remove drink called %v because of error: %v", name, err.Error())
		return
	} else {
		log.Infof("Removed drink called %v", name)
		c.JSON(http.StatusOK, OK)
	}
}
