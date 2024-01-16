FROM golang:1.21.5-alpine AS build

WORKDIR /go/src/github.com/kOstovic/CoffeeMachine
COPY . .
ENV CGO_ENABLED=0
RUN go mod download
RUN go build -o /bin/coffeemachine cmd/CoffeeMachineRestApiV3/main.go 

FROM golang:1.21.5-alpine
COPY --from=build /bin/coffeemachine /bin/coffeemachine
COPY --from=build /go/src/github.com/kOstovic/CoffeeMachine/cmd/CoffeeMachineRestApiV3/config.yaml /bin/config.yaml
ENV GIN_MODE=release
ENTRYPOINT ["/bin/coffeemachine"]