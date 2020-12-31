package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kOstovic/CoffeeMachine/internal/models"
)

type cofeeMachineController struct {
	Ingredients models.Ingredient
	Money       models.Denomination
}

//machineInitialized is private variable used for checking whether machine has been initialized
var (
	machineInitialized bool = false
)

func newCoffeeMachineController() *cofeeMachineController {
	return &cofeeMachineController{}
}

func (cmContoller cofeeMachineController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Hello from user controller"))
	if strings.ToLower(r.URL.Path) == "/coffeemachine" || strings.ToLower(r.URL.Path) == "/coffeemachine/" {
		switch r.Method {
		case http.MethodPost:
			cmContoller.postInitializeMachine(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

// InitializeMachine godoc
// @Summary Initialize Machine based on given Ingredient and money json
// @Produce json
// @Success 200 {object} cofeeMachineController
// @Router /coffeemachine [post]
func (cmContoller *cofeeMachineController) postInitializeMachine(w http.ResponseWriter, r *http.Request) {
	if machineInitialized == true {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Machine already Initialized"))
		return
	}
	iModel, mModel, err := cmContoller.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse Coffee Machine object\n" + err.Error()))
		return
	}
	cm, errIng := models.InitializeIngredients(iModel)
	if errIng != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not Initialize Coffee Machine object\n" + errIng.Error()))
		return
	}
	mm, errDen := models.InitializeDenominations(mModel)
	if errDen != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not Initialize Coffee Machine object\n" + errDen.Error()))
		return
	}
	machineInitialized = true
	encodeResponseAsJSON(fmt.Sprintf("%#v %#v", cm, mm), w)
}

func (cmContoller *cofeeMachineController) parseRequest(r *http.Request) (models.Ingredient, models.Denomination, error) {
	dec := json.NewDecoder(r.Body)
	var ing cofeeMachineController
	err := dec.Decode(&ing)

	if err != nil {
		return models.Ingredient{}, models.Denomination{}, err
	}

	if ing.Ingredients.Water < 0 || ing.Ingredients.Milk < 0 || ing.Ingredients.Sugar < 0 ||
		ing.Ingredients.CoffeeBeans < 0 || ing.Ingredients.TeaBeans < 0 || ing.Ingredients.Cups < 0 ||
		ing.Money.Half < 0 || ing.Money.One < 0 || ing.Money.Two < 0 || ing.Money.Five < 0 || ing.Money.Ten < 0 {
		return ing.Ingredients, ing.Money, fmt.Errorf("Values in ingredient and money cannot be negative'%v'", ing)
	}
	return ing.Ingredients, ing.Money, nil
}
