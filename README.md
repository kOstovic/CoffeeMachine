# CoffeeMachine
CoffeeMachine Rest API and Console APP implementation in GoLang used for learning GoLang and some devops techniques.

# Endpoints and Commands

## Endpoints V3 REST API version CoffeeMachine has endpoint:

Administrator endpoints are protected with Bearer token that you can get with /coffeemachine/login endpoint 
with login from pre-seeded account at startup from config.yaml

CoffeeMachine has endpoint (L means locked behind bearer token):
>
> - to get bearer token used in authorized only endpoints
>  - /coffeemachine/login
> - for initializing
>  - /coffeemachine L
> - for checking status of ingredient and updating Ingredient model
>  - /coffeemachine/ingredient L
> - for checking status of money and updating money based on Denomination model
>  - /coffeemachine/money L
> - for checking all available drinks, adding them and consuming them
>  - /coffeemachine/drinks
> - utility
>  - /metrics
> - for statistics of application
>  - /coffeemachine/statistics L
> - for statistics of coffeemachine
>  - /coffeemachine/health
> - for liveness probe
> -  /coffeemachine/swagger/index.html swagger endpoint in restapiv3

## Console APP version CoffeeMachine has command structure:

>
> - for initializing 
>  - coffeemachine subcommand
> - for checking status of ingredient and updating Ingredient model
>  - ingredient subcommand
> - for checking status of money and updating money based on Denomination model    
>  - money subcommand
> - for checking all available drinks, adding them and consuming them   
>  - drinks subcommand


# Running CoffeeMachineRestApi

## Config structure 
example for docker:
```yaml
database:
  type: "postgresql"
  host: "172.17.0.2"
  user: "postgres"
  password: "password"
  port: "5432"
  parameters: "sslmode=disable TimeZone=Europe/Zagreb"
  dbname_ingredient: "ingredient"
  dbname_denomination: "denomination"
  dbname_drinks: "drinks"
  initialized: "false"
log:
  level: "debug"
auth:
  username: "admin"
  password: "mypass"
```

## Docker

Dockerfile is in deployments/apiVersion/Dockerfile
Run BuildDocker.ps1 in api folder to automatically build docker image - if you previously set: "Set-ExecutionPolicy unrestricted" or just run BuildDocker.bat

After that you can run this docker image for example on some other port like this:
```sh
docker run -p 3000:3000 github.com/kostovic/coffeemachine/restapiv3:0.10.0 --env auth_password=newpass
```
restapiv3.0 is running in production mode by default

## Docker Compose

Docker Compose file is set in deployment/DockerCompose
Edit .env file to change variables from docker-compose file and then to run in detached mode run first command and to stop second command:
```sh
docker-compose -f docker-compose.full.yml up -d
docker-compose -f docker-compose.full.yml down
```
More in README.md in DockerCompose folder

## HELM chart

CoffeeMachineChart is Kubernetes set of yamls with own rules how to deploy to local or external cluster. Edit HELM/CoffeeMachineChart/values.yaml for more options and run or delete release with:

```sh
helm install my-release .
helm delete my-release
```

More in README.md in HELM/CoffeeMachineChart folder

## Manually 
Same is for CoffeeMachineRestApiV1 and CoffeeMachineRestApiV3
Run command in cmd/CoffeeMachineRestApi
```sh
go build -o coffeeMachine.exe main.go
```
and run coffeeMachine.exe

API is running in localhost:3000
(port will be customizable in future refactor)

Postman collection can be used for testing or just go to swagger endpoint in CoffeeMachineRestApiV3.
Metrics endpoint is exposed on /metrics 
Health endpoint is exposed on /coffeemachine/health
Logging level can be set with environment variable "LOG_LEVEL" in runtime

# Running CoffeeMachineConsole

## Manually

Run command in cmd/CoffeeMachineConsoleV1 or cmd/CoffeeMachineConsoleV3
```
go build -o coffeeMachine.exe main.go
```
and run coffeeMachine.exe

APP is running in local a console, commands are shown in console (or from help command in v3)
and example jsons for commands can be taken from postman collection


# Models

Models for V3 - DBModels have additional info and are backwards compatible.
Older models are still used because of in-memory part of V3 API and older models are for frontend REST API itself.
DB models uses additional GORM fields and have tenantName and Name field.

Models have tests inside their folder/package which can be run with command:
go test internal\models

```go
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

# Metrics

With low usage max MEM usage is around 30Mi and up to 1% CPU. GC doesnt seem to be active on low usage.

# License

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details or check LICENSE 
file in root of this repository.