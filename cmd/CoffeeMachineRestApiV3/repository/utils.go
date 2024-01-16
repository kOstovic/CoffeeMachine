package repository

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/kOstovic/CoffeeMachine/cmd/CoffeeMachineRestApiV3/config"
	"github.com/kOstovic/CoffeeMachine/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DbIngredient   *gorm.DB
	DbDenomination *gorm.DB
	DbDrinks       *gorm.DB
)

func InitDatabaseFromConfig() {
	initializedStartup()
	var wg sync.WaitGroup
	wg.Add(3)
	go initIngredientDatabase(&wg)
	go initDenominationDatabase(&wg)
	go initDrinksDatabase(&wg)
	wg.Wait()
}

func initializedStartup() {

	boolValue, err := strconv.ParseBool(config.Configuration.Database.INITIALIZED)
	if err != nil {
		panic(err)
	}
	if boolValue {
		MachineInitialized = boolValue
	}
}

func initIngredientDatabase(wg *sync.WaitGroup) {
	DbIngredient, _ = initDatabase(config.Configuration.Database.DBNAME_INGREDIENT)

	DbIngredient.AutoMigrate(&models.IngredientDB{})
	log.Debug("Ingredient Database initialized")
	defer wg.Done()
}

func initDenominationDatabase(wg *sync.WaitGroup) {
	DbDenomination, _ = initDatabase(config.Configuration.Database.DBNAME_DENOMINATION)

	log.Debug("Denomination Database initialized")
	DbDenomination.AutoMigrate(&models.DenominationDB{})
	defer wg.Done()
}

func initDrinksDatabase(wg *sync.WaitGroup) {
	DbDrinks, _ = initDatabase(config.Configuration.Database.DBNAME_DRINKS)

	DbDrinks.AutoMigrate(&models.DrinkDB{})
	log.Debug("Drinks Database initialized")
	defer wg.Done()
}

func initDatabase(dbname string) (*gorm.DB, error) {
	var dbObj *gorm.DB
	switch config.Configuration.Database.DB_TYPE {
	case "postgresql":
		dbObj, _ = initDatabasePostgres(dbname)
	default:
		err := fmt.Sprintf("Currenly postgresql is only supported database, your DB_TYPE: %v", config.Configuration.Database.DB_TYPE)
		log.Fatal(err)
		panic(err)
	}
	return dbObj, nil
}

func initDatabasePostgres(dbname string) (*gorm.DB, error) {
	var err error
	var dbOpen *gorm.DB
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s", config.Configuration.Database.DB_HOST, config.Configuration.Database.DB_PORT, config.Configuration.Database.DB_USER, config.Configuration.Database.DB_PASSWORD, dbname, config.Configuration.Database.DB_PARAMETERS)
	dbOpen, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err == nil && config.Configuration.Database.INITIALIZED == "false" {
		errString := fmt.Sprintf("Database %s opened and exist even thought it should not exist as INITIALIZED=%s", dbname, config.Configuration.Database.INITIALIZED)
		log.Fatal(errString)
		panic(errString)
	}
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") && config.Configuration.Database.INITIALIZED == "false" {
			dbOpen, _ = createDatabasePostgres(dbname)
		} else {
			errString := fmt.Sprintf("Database %s does not exist even thought it should exist as INITIALIZED=%s "+err.Error(), dbname, config.Configuration.Database.INITIALIZED)
			log.Fatal(errString)
			panic(errString)
		}
	}
	return dbOpen, nil
}

func createDatabasePostgres(dbname string) (*gorm.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s %s", config.Configuration.Database.DB_HOST, config.Configuration.Database.DB_PORT, config.Configuration.Database.DB_USER, config.Configuration.Database.DB_PASSWORD, config.Configuration.Database.DB_PARAMETERS)
	DB, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Could not connect to database: " + err.Error())
		panic(err)
	}
	log.Debugf("Creating database %v", dbname)
	createDatabaseCommand := fmt.Sprintf("CREATE DATABASE %v", dbname)
	createDbExec := DB.Exec(createDatabaseCommand)
	log.Debugf("Created database %v", dbname)
	if createDbExec.Error != nil {
		log.Fatalf("Could not create database: %s"+err.Error(), dbname)
		panic(createDbExec.Error)
	}

	connectionString = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s", config.Configuration.Database.DB_HOST, config.Configuration.Database.DB_PORT, config.Configuration.Database.DB_USER, config.Configuration.Database.DB_PASSWORD, dbname, config.Configuration.Database.DB_PARAMETERS)
	dbObj, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Could not connect to newly created database %v: "+err.Error(), dbname)
		panic(err)
	}
	return dbObj, nil
}
