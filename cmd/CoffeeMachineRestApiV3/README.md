<CoffeeMachine PoC software in golang>
Copyright (C) <2024>  <Krešimir Ostović>

# CoffeeMachine
CoffeeMachine Rest API implementation in GoLang used for learning GoLang.

## Endpoints
Administrator endpoints are protected with Bearer token that you can get with /coffeemachine/login endpoint 
with login from pre-seeded account at startup from config.yaml

CoffeeMachine has endpoint:
>
> - to get bearer token used in authorized only endpoints
>  - /coffeemachine/login
> - for initializing
>  - /coffeemachine
> - for checking status of ingredient and updating Ingredient model
>  - /coffeemachine/ingredient
> - for checking status of money and updating money based on Denomination model
>  - /coffeemachine/money
> - for checking all available drinks, adding them and consuming them
>  - /coffeemachine/drinks
> - utility
>  - /metrics
>  - /coffeemachine/health
> -  /coffeemachine/swagger/index.html swagger endpoint in restapiv3

## Running CoffeeMachineRestApiV3

While V2 was inmemory only, V3 uses database with implemenation supported currently in postgresql.

Run command in cmd/CoffeeMachineRestApiV3
```
go build -o coffeeMachine.exe main.go
```
and run coffeeMachine.exe

API is running in localhost:3000 
(port will be customizable in future refactor)

Postman collection can be used for testing.
Metrics endpoint is exposed on /metrics 
Health endpoint is exposed on /coffeemachine/health
Config parameters can be set either in config.yaml or as environment variables

config structure:
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

Currently, building Dockerfile from deployments/apiVersion/Dockerfile

Run BuildDocker.ps1 in api folder to automatically build docker image - if you previously set: Set-ExecutionPolicy unrestricted

Or just run BuildDocker.bat

After that you can run this docker image for example on some other port like this:
```
docker run -p 3000:3000 github.com/kostovic/coffeemachine/restapiv3:0.10.0 -d
```
