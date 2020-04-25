package models

import "fmt"

type Drink struct {
	Water       int
	Milk        int
	Sugar       int
	CoffeeBeans int
	TeaBeans    int
	Cups        int
	Money       float64
}

var (
	drinks = make(map[string]Drink)
)

func GetAvailableDrinks() map[string]Drink {
	return drinks
}

func AddDrink(name string, drink Drink) (map[string]Drink, error) {
	_, ok := drinks[name]
	if ok == true {
		return drinks, fmt.Errorf("Drink already exists '%v'", drinks[name])
	}
	drinks[name] = drink
	return drinks, nil
}

func ConsumeDrink(name string) (bool, error) {
	machineIngredients.Water -= drinks[name].Water
	machineIngredients.Milk -= drinks[name].Milk
	machineIngredients.Sugar -= drinks[name].Sugar
	machineIngredients.CoffeeBeans -= drinks[name].CoffeeBeans
	machineIngredients.TeaBeans -= drinks[name].TeaBeans
	machineIngredients.Cups -= drinks[name].Cups
	return true, nil
}

func CheckPrereqForDrink(name string) (bool, float64, error) {
	_, ok := drinks[name]
	if ok == true {
		if (machineIngredients.Water-drinks[name].Water) < 0 ||
			(machineIngredients.Milk-drinks[name].Milk) < 0 ||
			(machineIngredients.Sugar-drinks[name].Sugar) < 0 ||
			(machineIngredients.CoffeeBeans-drinks[name].CoffeeBeans) < 0 ||
			(machineIngredients.TeaBeans-drinks[name].TeaBeans) < 0 ||
			(machineIngredients.Cups-drinks[name].Cups) < 0 {
			return false, 0, fmt.Errorf("There are not enough ingredients for drink with name '%s'", name)
		}
		return true, drinks[name].Money, nil
	} else {
		return false, 0, fmt.Errorf("Drink with name '%s' not found", name)
	}
}
