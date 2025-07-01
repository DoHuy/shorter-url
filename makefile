APP_NAME=shorter-rest-api
BIN_DIR=bin
SWAGGER_DIR=docs

.PHONY: all build run test clean swag docker devbox

all: build

build:
    go build -o $(APP_NAME) .

run: build
    $(APP_NAME)

test:
    go test ./...


docker-build:
    docker build -t $(APP_NAME) .

docker-run:
    docker run -p 8080:8080 $(APP_NAME)

devbox: