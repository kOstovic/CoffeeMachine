#Copy over dependencies and build docker image
#Mock CI part

Copy-Item -Path "..\..\internal\" -Destination "CoffeeMachine/internal" -Recurse -Container
Copy-Item -Path "..\..\cmd\CoffeeMachineRestApiV2\" -Destination "CoffeeMachine\cmd\CoffeeMachineRestApiV2" -Recurse -Container
Copy-Item -Path "..\..\go.mod" -Destination "CoffeeMachine"
Copy-Item -Path "..\..\go.sum" -Destination "CoffeeMachine"
$VER=$(cat sem.ver) 
docker build -f Dockerfile -t github.com/kostovic/coffeemachine/restapiv2:$VER --no-cache "CoffeeMachine"
Remove-Item "CoffeeMachine" -Recurse