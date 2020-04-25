package controllers

import (
	"coffeeMachine/models"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type drinkController struct {
}

func newDrinkController() *drinkController {
	return &drinkController{}
}

func (dmContoller drinkController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Hello from user controller"))
	if strings.ToLower(r.URL.Path) == "/coffeemachine/drinks" || strings.ToLower(r.URL.Path) == "/coffeemachine/drinks/" {
		switch r.Method {
		case http.MethodGet:
			name, queryFound := r.URL.Query()["name"]
			if !queryFound || len(name[0]) < 1 {
				dmContoller.getAllAvailableDrinks(w, r)
			} else {
				dmContoller.getConsumeDrink(w, r)
			}
		case http.MethodPost:
			dmContoller.postAddDrink(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

// getAllAvailableDrinks godoc
// @Summary Get all drinks available
// @Produce json array
// @Success 200 {object} {string: models.drink}
// @Router /coffeemachine/drinks [get]
func (dmContoller *drinkController) getAllAvailableDrinks(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(models.GetAvailableDrinks(), w)
}

// ConsumeDrink godoc
// @Summary Consumes Dring over GET on /coffeemachine/drinks
// @Produce json
// @Success 200 {object} models.Denomination
// @Router /coffeemachine/drinks?name [get]
func (dmContoller *drinkController) getConsumeDrink(w http.ResponseWriter, r *http.Request) {
	name, queryFound := r.URL.Query()["name"]
	if !queryFound || len(name[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse name header request\n"))
		return
	}
	prereq, cost, err := models.CheckPrereqForDrink(name[0])
	if !prereq || err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	denominationParam := models.Denomination{Half: dmContoller.checkDenFromURL("half", r.URL),
		One: dmContoller.checkDenFromURL("one", r.URL), Two: dmContoller.checkDenFromURL("two", r.URL),
		Five: dmContoller.checkDenFromURL("five", r.URL), Ten: dmContoller.checkDenFromURL("ten", r.URL)}
	denRet, err := models.UpdateDenominationConsume(denominationParam, cost)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		strRet := fmt.Sprintf("\n%#v", denRet)
		w.Write([]byte(err.Error() + strRet))
		return
	}
	models.ConsumeDrink(name[0])
	encodeResponseAsJSON(denRet, w)
}

// AddDrink godoc
// @Summary Initialize new drink to consume on given Drink json
// @Produce json
// @Success 200 {object} model.Drink
// @Router /coffeemachine/drinks?name= [post]
func (dmContoller *drinkController) postAddDrink(w http.ResponseWriter, r *http.Request) {
	name, queryFound := r.URL.Query()["name"]
	if !queryFound || len(name[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Name for the drink should be provided in name query"))
		return
	}
	cm, err := dmContoller.parseRequestDrink(r)
	if err != nil || !queryFound {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse consumable object or name header\n" + err.Error()))
		return
	}
	drink, err := models.AddDrink(name[0], cm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(drink, w)
}

func (dmContoller *drinkController) checkDenFromURL(den string, u *url.URL) int {
	valueStr, queryFound := u.Query()[den]
	if !queryFound || len(valueStr[0]) < 1 {
		return 0
	}
	valueInt, err := strconv.Atoi(valueStr[0])
	if err != nil || valueInt < 0 {
		return 0
	}
	return valueInt
}

func (dmContoller *drinkController) parseRequestDrink(r *http.Request) (models.Drink, error) {
	dec := json.NewDecoder(r.Body)
	var d models.Drink
	err := dec.Decode(&d)
	if err != nil {
		return models.Drink{}, err
	}
	if d.Water < 0 || d.Milk < 0 || d.Sugar < 0 ||
		d.CoffeeBeans < 0 || d.TeaBeans < 0 || d.Cups < 0 ||
		d.Money < 0 {
		return d, fmt.Errorf("Values in drink cannot be negative'%v'", d)
	}
	return d, nil
}
