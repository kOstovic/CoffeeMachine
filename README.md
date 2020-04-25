# CoffeeMachine
CoffeeMachine Rest API implementation in GoLang used for learning GoLang.

## Endpoints
CoffeeMachine has endpoint:
>
> - for initializing
	/coffeemachine 
> - for checking status of ingredient and updating Ingredient model
	/coffeemachine/ingredient 
> - for checking status of money and updating money based on Denomination model    
    /coffeemachine/money 
> - for checking all available drinks, adding them and consuming them   
	/coffeemachine/drinks 


## Models
```
//model used for initializing machine
cofeeMachineController struct {
	Ingredients models.Ingredient
	Money       models.Denomination
}

//model used for adding Ingredient to the machine and internally consuming it
type Ingredient struct {
	Water       int
	Milk        int
	Sugar       int
	CoffeeBeans int
	TeaBeans    int
	Cups        int
}

//model used for adding drinks
Drink struct {
	Water       int
	Milk        int
	Sugar       int
	CoffeeBeans int
	TeaBeans    int
	Cups        int
	Money       float64
}

//model used for updating money and consuming drink
Denomination struct {
	Half  int
	One   int
	Two   int
	Five  int
	Ten   int
	Total float64
}
```
## Running

Run command 
```
go build coffeeMachine
```
and run coffeeMachine.exe

API is running in localhost:3000 
(port will be customizable in future refactor)

Postman collection can be used for testing.