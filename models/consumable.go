package models

import (
	"errors"
	"fmt"
	"reflect"
)

type Consumable struct {
	Water       int
	Milk        int
	Sugar       int
	CoffeeBeans int
	TeaBeans    int
	Cups        int
	Money       float32
}

var (
	drinks                         = make(map[string]Consumable)
	machineConsumables *Consumable = new(Consumable)
	machineInitialized bool        = false
)

func GetMachineConsumables() *Consumable {
	return machineConsumables
}

func InitializeConsumables(c Consumable) (Consumable, error) {
	if machineInitialized == true {
		return *machineConsumables, errors.New("Machine already initialized")
	}
	if c.Water <= 0 || c.Milk <= 0 || c.Sugar <= 0 ||
		c.CoffeeBeans <= 0 || c.TeaBeans <= 0 || c.Cups <= 0 ||
		c.Money <= 0 {
		return Consumable{}, errors.New("Initializing CoffeMachine must have some consumables to work")
	}
	machineConsumables = &c
	machineInitialized = true
	return *machineConsumables, nil
}

func GetIngredienteValueByName(ingredient string) (string, error) {
	r := reflect.ValueOf(*machineConsumables)
	for i := 0; i < r.NumField(); i++ {
		if ingredient == r.Type().Field(i).Name {
			return fmt.Sprintf("Field: %s Value: %v", r.Type().Field(i).Name, r.Field(i).Interface()), nil
		}
	}
	return "", fmt.Errorf("Ingredient with name '%s' not found", ingredient)
}

func UpdateConsumablePatch(c Consumable) (Consumable, error) {
	if c.Money > 0 {
		machineConsumables.Money = c.Money
	} else if c.Money < 0 {
		return *machineConsumables, fmt.Errorf("Ingredient with name money cannot be negative'%v'", c.Money)
	}
	if c.Water > 0 {
		machineConsumables.Water = c.Water
	}
	if c.Milk > 0 {
		machineConsumables.Milk = c.Milk
	}
	if c.Sugar > 0 {
		machineConsumables.Sugar = c.Sugar
	}
	if c.CoffeeBeans > 0 {
		machineConsumables.CoffeeBeans = c.CoffeeBeans
	}
	if c.TeaBeans > 0 {
		machineConsumables.TeaBeans = c.TeaBeans
	}
	if c.Cups > 0 {
		machineConsumables.Cups = c.Cups
	}
	return *machineConsumables, nil
}

func UpdateConsumablePut(c Consumable) (Consumable, error) {
	if c.Water < 0 || c.Milk < 0 || c.Sugar < 0 ||
		c.CoffeeBeans < 0 || c.TeaBeans < 0 || c.Cups < 0 ||
		c.Money < 0 {
		return Consumable{}, fmt.Errorf("Values in consumable cannot be negative'%v'", c)
	}
	machineConsumables.Money = c.Money
	machineConsumables.Water = c.Water
	machineConsumables.Milk = c.Milk
	machineConsumables.Sugar = c.Sugar
	machineConsumables.CoffeeBeans = c.CoffeeBeans
	machineConsumables.TeaBeans = c.TeaBeans
	machineConsumables.Cups = c.Cups

	return *machineConsumables, nil
}

func UpdateIngredienteValueByName(ingredient string, value int) (Consumable, error) {
	if value < 0 {
		return Consumable{}, fmt.Errorf("Value cannot be negative'%v'", value)
	}
	switch ingredient {
	case "Water":
		machineConsumables.Water = value
	case "Milk":
		machineConsumables.Milk = value
	case "Sugar":
		machineConsumables.Sugar = value
	case "CoffeeBeans":
		machineConsumables.CoffeeBeans = value
	case "TeaBeans":
		machineConsumables.TeaBeans = value
	case "Cups":
		machineConsumables.Cups = value
	default:
		return Consumable{}, fmt.Errorf("Ingredient with name '%s' not found", ingredient)
	}
	return *machineConsumables, nil
}

func UpdateMoney(value float32) (Consumable, error) {
	if value < 0 {
		return Consumable{}, fmt.Errorf("Money cannot be negative'%v'", value)
	} else {
		machineConsumables.Money = value
		return *machineConsumables, nil
	}
}

func GetAvailableDrinks() map[string]Consumable {
	return drinks
}

func AddDrink(name string, ingredients Consumable) (map[string]Consumable, error) {
	_, ok := drinks[name]
	if ok == true {
		return drinks, fmt.Errorf("Drink already exists '%v'", drinks[name])
	}
	drinks[name] = ingredients
	return drinks, nil
}

func ConsumeDrink(name string) (bool, error) {
	_, ok := drinks[name]
	if ok == true {
		if (machineConsumables.Water-drinks[name].Water) < 0 ||
			(machineConsumables.Milk-drinks[name].Milk) < 0 ||
			(machineConsumables.Sugar-drinks[name].Sugar) < 0 ||
			(machineConsumables.CoffeeBeans-drinks[name].CoffeeBeans) < 0 ||
			(machineConsumables.TeaBeans-drinks[name].TeaBeans) < 0 ||
			(machineConsumables.Cups-drinks[name].Cups) < 0 {
			return false, fmt.Errorf("There are not enough ingredients for drink with name '%s'", name)
		}

		machineConsumables.Water -= drinks[name].Water
		machineConsumables.Milk -= drinks[name].Milk
		machineConsumables.Sugar -= drinks[name].Sugar
		machineConsumables.CoffeeBeans -= drinks[name].CoffeeBeans
		machineConsumables.TeaBeans -= drinks[name].TeaBeans
		machineConsumables.Cups -= drinks[name].Cups
		machineConsumables.Money += drinks[name].Money
		return true, nil
	} else {
		return false, fmt.Errorf("Drink with name '%s' not found", name)
	}
}
