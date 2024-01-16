package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kOstovic/CoffeeMachine/internal/models"
	log "github.com/sirupsen/logrus"
)

func checkErrCode(err error) int {
	if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "already exists") {
		return http.StatusBadRequest
	} else if strings.Contains(err.Error(), "database") {
		return http.StatusInternalServerError
	} else {
		return http.StatusBadRequest
	}
}

func checkDenFromURL(c *gin.Context) (models.Denomination, error) {
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
		log.Debugf("Denomination could not be parsed from request %v", money)
		return models.Denomination{}, fmt.Errorf("Denomination could not be parsed")
	}
}

func checkDenFromBody(c *gin.Context) (models.Denomination, error) {
	var money Money
	if c.ShouldBindJSON(&money) == nil {
		log.Debugf("====== Bind By JSON ====== from request %v", money)
		return models.Denomination{Half: money.Half,
			One: money.One, Two: money.Two,
			Five: money.Five, Ten: money.Ten}, nil
	} else {
		log.Debugf("Denomination could not be parsed from request %v", money)
		return models.Denomination{}, fmt.Errorf("Denomination could not be parsed")
	}
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
		log.Debugf("Ingredient could not be parsed from request %v", ingredient)
		return models.Ingredient{}, fmt.Errorf("Ingredient could not be parsed")
	}
}

func checkIngredientsFromBody(c *gin.Context) (models.Ingredient, error) {
	var ingredient Ingredient
	if c.ShouldBindJSON(&ingredient) == nil {
		log.Debugf("====== Bind By JSON ====== from request %v", ingredient)
		return models.Ingredient{Water: ingredient.Water,
			Milk: ingredient.Milk, Sugar: ingredient.Sugar,
			CoffeeBeans: ingredient.CoffeeBeans, TeaBeans: ingredient.TeaBeans,
			Cups: ingredient.Cups}, nil
	} else {
		log.Debugf("Ingredient could not be parsed from request %v", ingredient)
		return models.Ingredient{}, fmt.Errorf("Ingredient could not be parsed")
	}
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
		log.Debugf("Drink could not be parsed from request %v", &drink)
		return models.Drink{}, fmt.Errorf("Drink could not be parsed in both query and body")
	}
}

func checkDrinkFromBody(c *gin.Context) (models.Drink, error) {
	var drink Drink
	if c.ShouldBindJSON(&drink) == nil {
		log.Debugf("====== Bind By JSON ====== from request %v", &drink)
		return models.Drink{Water: drink.Water,
			Milk: drink.Milk, Sugar: drink.Sugar,
			CoffeeBeans: drink.CoffeeBeans, TeaBeans: drink.TeaBeans,
			Cups: drink.Cups, Money: drink.Money}, nil
	} else {
		log.Debugf("Drink could not be parsed from request %v", &drink)
		return models.Drink{}, fmt.Errorf("Drink could not be parsed in both query and body")
	}
}
