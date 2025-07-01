APP_NAME=shorter-rest-api
BIN_DIR=bin
SWAGGER_DIR=docs

.PHONY: all build run test clean swag docker devbox

all: build

build:
    GOOS=linux GOARCH=amd64 go build -o  $(APP_NAME) .

run:
    ./$(APP_NAME)

test:
    go test ./...


docker-run:
    docker compose -f docker-compose.yml up -d