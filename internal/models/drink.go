package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Drink struct {
	Water       uint16
	Milk        uint16
	Sugar       uint16
	CoffeeBeans uint16
	TeaBeans    uint16
	Cups        uint16
	Money       float64
}

type DrinkDB struct {
	gorm.Model
	Name        string `gorm:"type:varchar(60);uniqueIndex:idx_name_tenantname"`
	TenantName  string `gorm:"type:varchar(60);uniqueIndex:idx_name_tenantname"`
	Water       uint16
	Milk        uint16
	Sugar       uint16
	CoffeeBeans uint16
	TeaBeans    uint16
	Cups        uint16
	Money       float64
}

var (
	drinks = make(map[string]Drink)
)

func (dbDrink *DrinkDB) ConvertDrinkDBToDrink() Drink {
	return Drink{
		Water:       dbDrink.Water,
		Milk:        dbDrink.Milk,
		Sugar:       dbDrink.Sugar,
		CoffeeBeans: dbDrink.CoffeeBeans,
		TeaBeans:    dbDrink.TeaBeans,
		Cups:        dbDrink.Cups,
		Money:       dbDrink.Money,
	}
}

func (drink *Drink) ConvertDrinkToDB(name string, tenantnName string) DrinkDB {
	return DrinkDB{
		Name:        name,
		TenantName:  tenantnName,
		Water:       drink.Water,
		Milk:        drink.Milk,
		Sugar:       drink.Sugar,
		CoffeeBeans: drink.CoffeeBeans,
		TeaBeans:    drink.TeaBeans,
		Cups:        drink.Cups,
		Money:       drink.Money,
	}
}

func (drink *Drink) ValidationDrink() (bool, error) {
	if drink.Money < 0 {
		return false, fmt.Errorf("Drink must have non negative values for ingredients and money %v", *drink)
	}
	if drink.Water == 0 && drink.Milk == 0 && drink.Sugar == 0 &&
		drink.CoffeeBeans == 0 && drink.TeaBeans == 0 && drink.Cups == 0 {
		return false, fmt.Errorf("Drink must have at least one 0 or positive value for ingridients %v", *drink)
	}
	return true, nil
}

func GetAvailableDrinks() map[string]Drink {
	return drinks
}

func GetAvailableDrinksName() []string {
	drinksNameArr := make([]string, 0, len(drinks))
	for key := range drinks {
		drinksNameArr = append(drinksNameArr, key)
	}
	return drinksNameArr
}

func GetDrinkByName(name string) Drink {
	return drinks[name]
}

func AddDrink(name string, drink Drink) (Drink, error) {
	validation, err := drink.ValidationDrink()
	if !validation {
		return drink, err
	}
	_, drinkExist := drinks[name]
	if drinkExist {
		return drink, fmt.Errorf("Drink already exists '%v'", drinks[name])
	}
	drinks[name] = drink
	return drink, nil
}

func RemoveDrink(name string) (bool, error) {
	_, drinkExist := drinks[name]
	if !drinkExist {
		return false, fmt.Errorf("Drink doesn't exists '%v'", drinks[name])
	} else {
		delete(drinks, name)
		return true, nil
	}
}

func ConsumeDrink(name string) (bool, error) {
	prereq, _, err := CheckPrereqForDrink(name)
	if prereq {
		machineIngredients.Water -= drinks[name].Water
		machineIngredients.Milk -= drinks[name].Milk
		machineIngredients.Sugar -= drinks[name].Sugar
		machineIngredients.CoffeeBeans -= drinks[name].CoffeeBeans
		machineIngredients.TeaBeans -= drinks[name].TeaBeans
		machineIngredients.Cups -= drinks[name].Cups
		return true, nil
	} else {
		return false, err
	}
}

func (dbDrink *DrinkDB) ConsumeDrinkDB() (bool, error) {
	prereq, _, newIng, err := dbDrink.CheckResourcePrereqForDrinkDB()
	if prereq {
		machineIngredients.Water = newIng.Water
		machineIngredients.Milk = newIng.Milk
		machineIngredients.Sugar = newIng.Sugar
		machineIngredients.CoffeeBeans = newIng.CoffeeBeans
		machineIngredients.TeaBeans = newIng.TeaBeans
		machineIngredients.Cups = newIng.Cups
		return true, nil
	} else {
		return false, err
	}
}

func (dbDrink *DrinkDB) CheckResourcePrereqForDrinkDB() (bool, float64, Ingredient, error) {
	if (int(machineIngredients.Water)-int(dbDrink.Water)) < 0 ||
		(int(machineIngredients.Milk)-int(dbDrink.Milk)) < 0 ||
		(int(machineIngredients.Sugar)-int(dbDrink.Sugar)) < 0 ||
		(int(machineIngredients.CoffeeBeans)-int(dbDrink.CoffeeBeans)) < 0 ||
		(int(machineIngredients.TeaBeans)-int(dbDrink.TeaBeans)) < 0 ||
		(int(machineIngredients.Cups)-int(dbDrink.Cups)) < 0 {
		return false, 0, Ingredient{}, fmt.Errorf("There are not enough ingredients for drink with name '%s'", dbDrink.Name)
	}
	return true, dbDrink.Money,
		Ingredient{Water: uint16(machineIngredients.Water) - uint16(dbDrink.Water),
			Milk:        uint16(machineIngredients.Milk) - uint16(dbDrink.Milk),
			Sugar:       uint16(machineIngredients.Sugar) - uint16(dbDrink.Sugar),
			CoffeeBeans: uint16(machineIngredients.CoffeeBeans) - uint16(dbDrink.CoffeeBeans),
			TeaBeans:    uint16(machineIngredients.TeaBeans) - uint16(dbDrink.TeaBeans),
			Cups:        uint16(machineIngredients.Cups) - uint16(dbDrink.Cups)}, nil
}

func CheckPrereqForDrink(name string) (bool, float64, error) {
	_, drinkExist := drinks[name]
	if drinkExist {
		if (int(machineIngredients.Water)-int(drinks[name].Water)) < 0 ||
			(int(machineIngredients.Milk)-int(drinks[name].Milk)) < 0 ||
			(int(machineIngredients.Sugar)-int(drinks[name].Sugar)) < 0 ||
			(int(machineIngredients.CoffeeBeans)-int(drinks[name].CoffeeBeans)) < 0 ||
			(int(machineIngredients.TeaBeans)-int(drinks[name].TeaBeans)) < 0 ||
			(int(machineIngredients.Cups)-int(drinks[name].Cups)) < 0 {
			return false, 0, fmt.Errorf("There are not enough ingredients for drink with name '%s'", name)
		}
		return true, drinks[name].Money, nil
	} else {
		return false, 0, fmt.Errorf("Drink with name '%s' not found", name)
	}
}
