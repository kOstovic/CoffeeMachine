package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/kOstovic/CoffeeMachine/cmd/CoffeeMachineRestApiV3/config"
	"github.com/kOstovic/CoffeeMachine/cmd/CoffeeMachineRestApiV3/controllers"
	_ "github.com/kOstovic/CoffeeMachine/cmd/CoffeeMachineRestApiV3/docs"
	"github.com/kOstovic/CoffeeMachine/cmd/CoffeeMachineRestApiV3/repository"
)

// @title CoffeeMachine Swagger API
// @version openapi: 3.0.0
// @description Swagger API for Golang Project CoffeeMachine.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @query.collection.format multi
// @license.name GNU AFFERO
// @license.url https://github.com/kOstovic/coffeemachine/blob/master/LICENSE
// @BasePath /coffeemachine
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// Package classification of Product API
func main() {

	if err := config.LoadConfig(&config.Configuration); err != nil {
		log.Fatal(err.Error())
		panic(err)
	}
	var lvl log.Level = log.DebugLevel
	var err error
	lvls, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl, err = log.ParseLevel(config.Configuration.Log.LOG_LEVEL)
		if err != nil {
			lvl = log.DebugLevel
		}
	} else {
		lvl, err = log.ParseLevel(lvls)
		if err != nil {
			lvl = log.DebugLevel
		}
	}
	// set global log level
	log.SetLevel(lvl)
	log.Warnln("<CoffeeMachine>  Copyright (C) <2024>  <Kresimir Ostovic> This program comes with ABSOLUTELY NO WARRANTY; This is free software, and you are welcome to redistribute itunder certain conditions of AGPL 3.0 License.; Check LICENSE file in root of github repo or find license here: https://www.gnu.org/licenses/agpl-3.0.en.html")

	repository.InitDatabaseFromConfig()

	router := gin.New()
	router.Use(gin.LoggerWithWriter(log.StandardLogger().WriterLevel(log.DebugLevel)))
	router.Use(gin.Recovery())

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	coffeemachine := router.Group("/coffeemachine")
	controllers.RegisterRoutesCoffeeMachine(coffeemachine.Group("/"))
	controllers.RegisterRoutesAuth(coffeemachine.Group("/login"))
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
