package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/kOstovic/CoffeeMachine/cmd/CoffeeMachineRestApiV2/controllers"
	_ "github.com/kOstovic/CoffeeMachine/cmd/CoffeeMachineRestApiV2/docs"
)

// @title CoffeeMachine Swagger API
// @version 2.0
// @description Swagger API for Golang Project CoffeeMachine.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @query.collection.format multi
// @license.name Apache 2.0
// @license.url https://github.com/kOstovic/coffeemachine/blob/master/LICENSE

// @BasePath /coffeemachine
// Package classification of Product API.
func main() {

	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "debug"
	}
	ll, err := log.ParseLevel(lvl)
	if err != nil {
		ll = log.DebugLevel
	}
	// set global log level
	log.SetLevel(ll)

	router := gin.New()
	router.Use(gin.LoggerWithWriter(log.StandardLogger().WriterLevel(log.DebugLevel)))
	router.Use(gin.Recovery())

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	coffeemachine := router.Group("/coffeemachine")
	controllers.RegisterRoutesCoffeeMachine(coffeemachine.Group("/"))
	controllers.RegisterRoutesDrink(coffeemachine.Group("/drinks"))
	controllers.RegisterRoutesIngredients(coffeemachine.Group("/ingredients"))
	controllers.RegisterRoutesDenomination(coffeemachine.Group("/money"))
	router.GET("/coffeemachine/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//health test
	coffeemachine.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	//router := controllers.SetupRouter()
	// Listen and Server in 0.0.0.0:3000
	router.Run(":3000")
}
