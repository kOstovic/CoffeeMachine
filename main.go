package main

import (
	controllers "coffeeMachine/controller"
	"net/http"
)

func main() {
	controllers.RegisterControllers()
	http.ListenAndServe(":3000", nil)
}
