# CoffeeMachine
CoffeeMachine Rest API implementation in GoLang used for learning GoLang.

## Endpoints
CoffeeMachine has endpoint:
>
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
> -  /coffeemachine/swagger/index.html swagger endpoint in restAPIv2

## Running CoffeeMachineRestApiV2

Implementation is in-memory only.

Run command in cmd/CoffeeMachineRestApiV2
```
go build -o coffeeMachine.exe main.go
```
and run coffeeMachine.exe

API is running in localhost:3000 
(port will be customizable in future refactor)

Postman collection can be used for testing.
Metrics endpoint is exposed on /metrics 
Health endpoint is exposed on /coffeemachine/health
Logging level can be set with environment variable "LOG_LEVEL" in runtime

## Docker

Currently, building Dockerfile from deployments/apiVersion/Dockerfile

Run BuildDocker.ps1 in api folder to automatically build docker image - if you previously set: Set-ExecutionPolicy unrestricted

Or just run BuildDocker.bat

After that you can run this docker image for example on some other port like this:
```
docker run -p 3002:3000 github.com/kostovic/coffeemachine/restapiv2:0.9.0
```
