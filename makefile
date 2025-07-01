APP_NAME=shorter-rest-api
BIN_DIR=bin
SWAGGER_DIR=docs

.PHONY: all build run test clean swag docker devbox

all: build

build:
    go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/...

run: build
    ./$(BIN_DIR)/$(APP_NAME)

test:
    go test ./...

clean:
    rm -rf $(BIN_DIR)/
    rm -rf $(SWAGGER_DIR)/

swag:
    swag init --output $(SWAGGER_DIR)

docker-build:
    docker build -t $(APP_NAME) .

docker-run:
    docker run -p 8080:8080 $(APP_NAME)

devbox: