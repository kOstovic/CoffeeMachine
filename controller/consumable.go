package controllers

import (
	"coffeShop/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type coffeeMachineController struct {
}

func newCoffeeMachineController() *coffeeMachineController {
	return &coffeeMachineController{}
}

func (cmContoller coffeeMachineController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Hello from user controller"))
	if r.URL.Path == "/coffeemachine/ingredients" {
		switch r.Method {
		case http.MethodGet:
			name := r.Header.Get("ingredient")
			if name == "" {
				cmContoller.getAllIngredients(w, r)
			} else {
				cmContoller.getIngredientsByName(w, r)
			}
		case http.MethodPost:
			cmContoller.postInitializeMachine(w, r)
		case http.MethodPut:
			name := r.Header.Get("ingredient")
			if name == "" {
				cmContoller.putIngredients(w, r)
			} else {
				cmContoller.putIngredientsByName(w, r)
			}

		case http.MethodPatch:
			cmContoller.patchIngredients(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	} else if r.URL.Path == "/coffeemachine/drinks" {
		switch r.Method {
		case http.MethodGet:
			name := r.Header.Get("name")
			if name == "" {
				cmContoller.getAllAvailableDrinks(w, r)
			} else {
				cmContoller.getConsumeDrink(w, r)
			}
		case http.MethodPost:
			cmContoller.postAddDrink(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

func (cmContoller *coffeeMachineController) getAllIngredients(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(models.GetMachineConsumables(), w)
}

func (cmContoller *coffeeMachineController) getIngredientsByName(w http.ResponseWriter, r *http.Request) {
	name := r.Header.Get("ingredient")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse name header request\n"))
		return
	}
	cm, err := models.GetIngredienteValueByName(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	encodeResponseAsJSON(cm, w)
}

func (cmContoller *coffeeMachineController) postInitializeMachine(w http.ResponseWriter, r *http.Request) {
	cm, err := cmContoller.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse Coffee Machine object\n" + err.Error()))
		return
	}
	cm, err = models.InitializeConsumables(cm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not Initialize Coffee Machine object\n" + err.Error()))
		return
	}
	encodeResponseAsJSON(cm, w)
}

func (cmContoller *coffeeMachineController) putIngredients(w http.ResponseWriter, r *http.Request) {
	cm, err := cmContoller.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse consumable object"))
		return
	}
	cm, err = models.UpdateConsumablePut(cm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(cm, w)
}

func (cmContoller *coffeeMachineController) putIngredientsByName(w http.ResponseWriter, r *http.Request) {
	name := r.Header.Get("ingredient")
	valueStr := r.Header.Get("value")
	if name == "" || valueStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse name or value header request\n"))
		return
	}
	if name == "Money" {
		value, err := strconv.ParseFloat(valueStr, 32)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		cm, err := models.UpdateMoney(float32(value))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		encodeResponseAsJSON(cm, w)
	} else {
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		cm, err := models.UpdateIngredienteValueByName(name, value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		encodeResponseAsJSON(cm, w)
	}
}

func (cmContoller *coffeeMachineController) patchIngredients(w http.ResponseWriter, r *http.Request) {
	cm, err := cmContoller.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse consumable object"))
		return
	}
	cm, err = models.UpdateConsumablePatch(cm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(cm, w)
}

func (cmContoller *coffeeMachineController) getAllAvailableDrinks(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(models.GetAvailableDrinks(), w)
}

func (cmContoller *coffeeMachineController) getConsumeDrink(w http.ResponseWriter, r *http.Request) {
	name := r.Header.Get("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse name header request\n"))
		return
	}
	cm, err := models.ConsumeDrink(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(cm, w)
}

func (cmContoller *coffeeMachineController) postAddDrink(w http.ResponseWriter, r *http.Request) {
	name := r.Header.Get("name")
	cm, err := cmContoller.parseRequest(r)
	if err != nil || name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse consumable object or name header\n" + err.Error()))
		return
	}
	drink, err := models.AddDrink(name, cm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(drink, w)
}

func (cmContoller *coffeeMachineController) parseRequest(r *http.Request) (models.Consumable, error) {
	dec := json.NewDecoder(r.Body)
	var c models.Consumable
	err := dec.Decode(&c)
	if err != nil {
		return models.Consumable{}, err
	}
	if c.Water < 0 || c.Milk < 0 || c.Sugar < 0 ||
		c.CoffeeBeans < 0 || c.TeaBeans < 0 || c.Cups < 0 ||
		c.Money < 0 {
		return c, fmt.Errorf("Values in consumable cannot be negative'%v'", c)
	}
	return c, nil
}
