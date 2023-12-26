// package controllers
package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/kOstovic/CoffeeMachine/internal/models"
	log "github.com/sirupsen/logrus"
)

type Denomination struct {
	Half int `form:"Half" json:"Half" binding:"required_without_all=One Two Five Ten,validateDenomination"`
	One  int `form:"One" json:"One" binding:"required_without_all=Half Two Five Ten,validateDenomination"`
	Two  int `form:"Two" json:"Two" binding:"required_without_all=One Half Five Ten,validateDenomination"`
	Five int `form:"Five" json:"Five" binding:"required_without_all=One Two Half Ten,validateDenomination"`
	Ten  int `form:"Ten" json:"Ten" binding:"required_without_all=One Two Five Half,validateDenomination"`
}

// used for initialization of CoffeeMachine
type CoffeeMachine struct {
	Ingredients  Ingredient   `form:"Ingredients" json:"Ingredients"`
	Denomination Denomination `form:"Money" json:"Money" binding:"validateDenomination"`
}

// machineInitialized is private variable used for checking whether machine has been initialized
var (
	machineInitialized bool = false
)

// register route for coffeemachine in gin framework
func RegisterRoutesCoffeeMachine(router *gin.RouterGroup) {
	router.POST("", postInitializeMachine)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validateDenomination", validateDenomination)
		v.RegisterStructValidation(validateIngredient, Ingredient{})
	}
}

// InitializeMachine godoc
// @Summary Initialize Machine based on given Ingredient and money json
// @Description Initialize Machine based on given Ingredient and money json
// @Param CoffeeMachine body CoffeeMachine true "init CoffeeMachine object"
// @Accept json
// @Produce json
// @Success 200 {object} CoffeeMachine
// @Failure 400,404
// @Failure 500
// @Router / [post]
func postInitializeMachine(c *gin.Context) {
	if machineInitialized == true {
		log.Errorf("coffeeMachine cannot be initialized more than once")
		c.JSON(http.StatusBadRequest, "Machine already Initialized")
		return
	}
	iModel, mModel, err := checkCoffeeMachineFromReq(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		log.Errorf("coffeeMachine could not be initialized " + err.Error())
		return
	}
	cm, errIng := models.InitializeIngredients(iModel)
	if errIng != nil {
		c.JSON(http.StatusBadRequest, "Could not Initialize Coffee Machine object "+err.Error())
		log.Errorf("coffeeMachine could not be initialized " + err.Error())
		return
	}
	mm, errDen := models.InitializeDenominations(mModel)
	if errDen != nil {
		c.JSON(http.StatusBadRequest, "Could not Initialize Coffee Machine object "+err.Error())
		log.Errorf("coffeeMachine could not be initialized " + err.Error())
		return
	}
	machineInitialized = true
	log.Infof("coffeeMachine initialized with following parameters: Ingredients: %v Money: %v", cm, mm)
	c.JSON(http.StatusOK, fmt.Sprintf("Ingredients: %v Money: %v", cm, mm))
}

func checkCoffeeMachineFromReq(c *gin.Context) (models.Ingredient, models.Denomination, error) {
	var coffeeMachine CoffeeMachine
	if c.ShouldBindJSON(&coffeeMachine) == nil {
		log.Debugf("====== Bind By JSON ====== from request %v", coffeeMachine)
		return models.Ingredient{Water: coffeeMachine.Ingredients.Water,
				Milk: coffeeMachine.Ingredients.Milk, Sugar: coffeeMachine.Ingredients.Sugar,
				CoffeeBeans: coffeeMachine.Ingredients.CoffeeBeans, TeaBeans: coffeeMachine.Ingredients.TeaBeans,
				Cups: coffeeMachine.Ingredients.Cups},
			models.Denomination{Half: coffeeMachine.Denomination.Half,
				One: coffeeMachine.Denomination.One, Two: coffeeMachine.Denomination.Two,
				Five: coffeeMachine.Denomination.Five, Ten: coffeeMachine.Denomination.Ten}, nil
	} else {
		log.Warnf("coffeeMachine could not be parsed from request %v", coffeeMachine)
		return models.Ingredient{}, models.Denomination{}, fmt.Errorf("coffeeMachine could not be parsed or validation failed - check your values again")
	}
}

var validateDenomination validator.Func = func(fl validator.FieldLevel) bool {
	den, ok := fl.Field().Interface().(int)
	if ok {
		if den < 0 {
			return false
		}
	}
	return true
}

func validateIngredient(sl validator.StructLevel) {
	ing := sl.Current().Interface().(Ingredient)
	if (ing.Water < 0 || ing.Milk < 0 || ing.Sugar < 0 ||
		ing.CoffeeBeans < 0 || ing.TeaBeans < 0 || ing.Cups < 0) ||
		(ing.Water <= 0 && ing.Milk <= 0 && ing.Sugar <= 0 &&
			ing.CoffeeBeans <= 0 && ing.TeaBeans <= 0 && ing.Cups <= 0) {
		sl.ReportError(ing, "One of ingredients is not valid", "Ingredient", "Ingredient", "")
	}

}
