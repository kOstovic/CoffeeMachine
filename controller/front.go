package controllers

import (
	"encoding/json"
	"io"
	"net/http"
)

func RegisterControllers() {
	cmContoller := newCoffeeMachineController()

	http.Handle("/coffeemachine/ingredients", *cmContoller)
	http.Handle("/coffeemachine/ingredients/", *cmContoller)
	http.Handle("/coffeemachine/drinks", *cmContoller)
	http.Handle("/coffeemachine/drinks/", *cmContoller)
}

func encodeResponseAsJSON(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}
