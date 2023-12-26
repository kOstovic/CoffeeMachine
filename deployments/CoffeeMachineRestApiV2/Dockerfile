FROM golang:1.21.5-alpine AS build

WORKDIR /go/src/github.com/kOstovic/CoffeeMachine
COPY . ..
ENV CGO_ENABLED=0
RUN go mod download
RUN pwd && ls -alh
RUN go build -o /bin/coffeemachine cmd/CoffeeMachineRestApiV2/main.go 

FROM scratch
COPY --from=build /bin/coffeemachine /bin/coffeemachine
ENV GIN_MODE=release
ENTRYPOINT ["/bin/coffeemachine"]