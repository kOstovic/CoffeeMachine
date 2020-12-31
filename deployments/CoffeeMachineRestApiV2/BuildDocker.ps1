#Copy over dependencies and build docker image
#Mock CI/CD part

Copy-Item -Path "..\..\internal\" -Destination "CoffeeMachine/internal" -Recurse -Container
Copy-Item -Path "..\..\cmd\CoffeeMachineRestApiV2\" -Destination "CoffeeMachine\cmd\CoffeeMachineRestApiV2" -Recurse -Container

docker build . -t github.com/kostovic/coffeemachine/restapiv2:0.1 --no-cache
Remove-Item "CoffeeMachine" -Recurse