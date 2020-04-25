package controllers

import (
	"coffeeMachine/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type moneyController struct {
}

func newMoneyController() *moneyController {
	return &moneyController{}
}

func (mmContoller moneyController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Hello from user controller"))
	if strings.ToLower(r.URL.Path) == "/coffeemachine/money" || strings.ToLower(r.URL.Path) == "/coffeemachine/money/" {
		switch r.Method {
		case http.MethodGet:
			name, queryFound := r.URL.Query()["name"]
			if !queryFound || len(name[0]) < 1 {
				mmContoller.getAllAvailableDenomination(w, r)
			} else {
				mmContoller.getDenominationByName(w, r)
			}
		case http.MethodPut:
			name, queryFound := r.URL.Query()["name"]
			if !queryFound || len(name[0]) < 1 {
				mmContoller.putDenomination(w, r)
			} else {
				mmContoller.putDenominationByName(w, r)
			}

		case http.MethodPatch:
			mmContoller.patchDenomination(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

// getAllAvailableDenomination godoc
// @Summary Get all denominations available
// @Produce json
// @Success 200 {object} models.Denominationm
// @Router /coffeemachine/money [get]
func (mmContoller *moneyController) getAllAvailableDenomination(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(models.GetCurrentMoney(), w)
}

// getDenominationByName godoc
// @Summary Get denominations by name from query
// @Produce json
// @Success 200 {object} models.Denomination.Field
// @Router /coffeemachine/money?name= [get]
func (mmContoller *moneyController) getDenominationByName(w http.ResponseWriter, r *http.Request) {
	//name := r.Header.Get("denomination")
	name, queryFound := r.URL.Query()["name"]
	if !queryFound || len(name[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse name header request\n"))
		return
	}
	cm, err := models.GetDenominationValueByName(name[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(cm, w)
}

// putDenomination godoc
// @Summary Update Denomination based on given Denomination json, updates all
// @Produce json
// @Success 200 {object} model.Denomination
// @Router /coffeemachine/denomination [put]
func (mmContoller *moneyController) putDenomination(w http.ResponseWriter, r *http.Request) {
	cm, err := mmContoller.parseRequestDenomination(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse consumable object"))
		return
	}
	cm, err = models.UpdateDenominationPut(cm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(cm, w)
}

// putDenominationByName godoc
// @Summary Update denomination based on given Denomination name and value in query
// @Produce json
// @Success 200 {object} model.Denomination
// @Router /coffeemachine/denomination?name=&value [put]
func (mmContoller *moneyController) putDenominationByName(w http.ResponseWriter, r *http.Request) {
	name, queryNameFound := r.URL.Query()["name"]
	valueStr, queryValueFound := r.URL.Query()["value"]
	if !queryNameFound || !queryValueFound || len(name[0]) < 1 || len(valueStr[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse name or value header request\n"))
		return
	}
	value, err := strconv.Atoi(valueStr[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	cm, err := models.UpdateDenominationValueByName(name[0], value)
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
func (mmContoller *moneyController) patchDenomination(w http.ResponseWriter, r *http.Request) {
	cm, err := mmContoller.parseRequestDenomination(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not parse denomination object"))
		return
	}
	cm, err = models.UpdateDenominationPatch(cm)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	encodeResponseAsJSON(cm, w)
}

func (mmContoller *moneyController) parseRequestDenomination(r *http.Request) (models.Denomination, error) {
	dec := json.NewDecoder(r.Body)
	var d models.Denomination
	err := dec.Decode(&d)
	if err != nil {
		return models.Denomination{}, err
	}
	if d.Half < 0 || d.One < 0 || d.Two < 0 ||
		d.Five < 0 || d.Ten < 0 {
		return d, fmt.Errorf("Values in Denomination cannot be negative'%v'", d)
	}
	return d, nil
}
