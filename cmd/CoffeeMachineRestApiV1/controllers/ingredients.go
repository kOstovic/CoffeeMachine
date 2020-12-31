package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/kOstovic/CoffeeMachine/internal/models"
)

type ingredientsController struct {
}

func newIngredientsController() *ingredientsController {
	return &ingredientsController{}
}

func (imContoller ingredientsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Hello from ingredients controller"))
	if strings.ToLower(r.URL.Path) == "/coffeemachine/ingredients" || strings.ToLower(r.URL.Path) == "/coffeemachine/ingredients/" {
		switch r.Method {
		case http.MethodGet:
			name, queryFound := r.URL.Query()["name"]
			//name := r.Header.Get("ingredient")
			if !queryFound || len(name[0]) < 1 {
				imContoller.getAllIngredients(w, r)

			} else {
				imContoller.getIngredientsByName(w, r)
			}
		case http.MethodPut:
			name, queryFound := r.URL.Query()["name"]
			//name := r.Header.Get("ingredient")
			if !queryFound || len(name[0]) < 1 {
				imContoller.putIngredients(w, r)
			} else {
				imContoller.putIngredientsByName(w, r)
			}

		case http.MethodPatch:
			imContoller.patchIngredients(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

// getAllIngredients godoc
// @Summary Get all ingredients available
// @Produce json
// @Success 200 {object} {string: models.Ingredients}
// @Router /coffeemachine/ingredients [get]
func (imContoller *ingredientsController) getAllIngredients(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(models.GetMachineIngredients(), w)
}

// getIngredientsByName godoc
// @Summary Get ingredient by name from query
// @Produce json
// @Success 200 {object} models.Ingredients.Field
// @Router /coffeemachine/ingredients?name= [get]
func (imContoller *ingredientsController) getIngredientsByName(w http.ResponseWriter, r *http.Request) {
	name, queryFound := r.URL.Query()["name"]
	if !queryFound || len(name[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse name header request\n"))
		return
	}
	cm, err := models.GetIngredienteValueByName(name[0])
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	encodeResponseAsJSON(cm, w)
}

// putIngredients godoc
// @Summary Update ingredients based on given Ingredient json, updates all
// @Produce json
// @Success 200 {object} model.Ingredient
// @Router /coffeemachine/ingredients [put]
func (imContoller *ingredientsController) putIngredients(w http.ResponseWriter, r *http.Request) {
	cm, err := imContoller.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse ingredient object"))
		return
	}
	cm, err = models.UpdateIngredientPut(cm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(cm, w)
}

// putIngredientsByName godoc
// @Summary Update ingredients based on given Ingredient name and value in query
// @Produce json
// @Success 200 {object} model.Ingredient
// @Router /coffeemachine/ingredients?name=&value [put]
func (imContoller *ingredientsController) putIngredientsByName(w http.ResponseWriter, r *http.Request) {
	name, queryNameFound := r.URL.Query()["name"]
	valueStr, queryValueFound := r.URL.Query()["value"]
	if !queryNameFound || !queryValueFound || len(name[0]) < 1 || len(valueStr[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse name or value header request\n"))
		return
	}
	value, err := strconv.ParseUint(valueStr[0], 10, 16)
	valueuint16 := uint16(value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	cm, err := models.UpdateIngredientValueByName(name[0], valueuint16)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(cm, w)
}

// patchIngredients godoc
// @Summary Update ingredients based on given Ingredient json, update only given
// @Produce json
// @Success 200 {object} model.Ingredient
// @Router /coffeemachine/ingredients [patch]
func (imContoller *ingredientsController) patchIngredients(w http.ResponseWriter, r *http.Request) {
	cm, err := imContoller.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse ingredient object"))
		return
	}
	cm, err = models.UpdateIngredientPatch(cm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(cm, w)
}

func (cmContoller *ingredientsController) parseRequest(r *http.Request) (models.Ingredient, error) {
	dec := json.NewDecoder(r.Body)
	var ing models.Ingredient
	err := dec.Decode(&ing)
	if err != nil {
		return models.Ingredient{}, err
	}
	if ing.Water < 0 || ing.Milk < 0 || ing.Sugar < 0 ||
		ing.CoffeeBeans < 0 || ing.TeaBeans < 0 || ing.Cups < 0 {
		return ing, fmt.Errorf("Values in ingredient cannot be negative'%v'", ing)
	}
	return ing, nil
}
