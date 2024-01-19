// package controllers
package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/kOstovic/CoffeeMachine/cmd/CoffeeMachineRestApiV3/repository"
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

// register route for coffeemachine in gin framework
func RegisterRoutesCoffeeMachine(router *gin.RouterGroup) {
	router.POST("", postInitializeMachine)
	router.DELETE("", deleteDeInitializeMachine)
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
// @Security BearerAuth
func postInitializeMachine(c *gin.Context) {
	if code, err := parseToken(c); err != nil {
		c.JSON(code, err.Error())
		return
	}
	if repository.MachineInitialized == true {
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
	_, errRepo := repository.InitializeMachine(iModel, mModel)
	if errRepo != nil {
		c.JSON(checkErrCode(errRepo), "Could not Initialize Coffee Machine object "+errRepo.Error())
		log.Errorf("coffeeMachine could not be initialized " + errRepo.Error())
		return
	}
	mModel.CalculateTotal()

	log.Infof("coffeeMachine initialized with following parameters: Ingredients: %v Money: %v", iModel, mModel)
	c.JSON(http.StatusOK, fmt.Sprintf("Ingredients: %v Money: %v", iModel, mModel))
}

// deleteDeInitializeMachine godoc
// @Summary DeInitialize Machine
// @Description DeInitialize Machine based
// @Accept json
// @Produce json
// @Success 200 {object} CoffeeMachine
// @Failure 400,404
// @Failure 500
// @Router / [delete]
// @Security BearerAuth
func deleteDeInitializeMachine(c *gin.Context) {
	if code, err := parseToken(c); err != nil {
		c.JSON(code, err.Error())
		return
	}
	if repository.MachineInitialized == false {
		log.Errorf("coffeeMachine cannot be deinitialized more than once")
		c.JSON(http.StatusBadRequest, "Machine already DeInitialized")
		return
	}
	_, errdeinit := repository.DeleteDeInitializeMachine()
	if errdeinit != nil {
		c.JSON(checkErrCode(errdeinit), "Could not DeInitialize Coffee Machine object "+errdeinit.Error())
		return
	}
	log.Infof("coffeeMachine deinitialized")
	c.JSON(http.StatusOK, fmt.Sprintf("coffeeMachine deinitialized"))
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
		log.Debugf("coffeeMachine could not be parsed from request %v", coffeeMachine)
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
