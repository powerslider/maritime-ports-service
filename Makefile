envfile ?= .env
-include $(envfile)
ifneq ("$(wildcard $(envfile))","")
	export $(shell sed 's/=.*//' $(envfile))
endif

GOLANGCI_VERSION:=1.52.2
SWAG_VERSION:=v1.8.12
PROJECT_NAME:=maritime-ports-service
GOPATH_BIN:=$(shell go env GOPATH)/bin

.PHONY: install
install:
	# Install Swag tool for Swagger API documentation generation.
	go install \
		github.com/swaggo/swag/cmd/swag@${SWAG_VERSION}

	# Install golangci-lint for go code linting.
	curl -sSfL \
		"https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" | \
		sh -s -- -b ${GOPATH_BIN} v${GOLANGCI_VERSION}

.PHONY: all
all: clean init lint test build-server

.PHONY: init
init:
	@cp .env.dist .env

.PHONY: lint
lint:
	@echo ">>> Performing golang code linting.."
	golangci-lint run --config=.golangci.yml

.PHONY: test
test:
	@echo ">>> Running Unit Tests..."
	go test -v -race ./...

.PHONY: cover-test
cover-test:
	@echo ">>> Running Tests with Coverage..."
	go test -v -race ./... -coverprofile=coverage.txt -covermode=atomic

.PHONY: build-server
build-server:
	@echo ">>> Building ${PROJECT_NAME} API server..."
	go build -o bin/${PROJECT_NAME} cmd/${PROJECT_NAME}/main.go

.PHONY: run-server
run-server:
	@echo ">>> Running ${PROJECT_NAME} API server..."
	@go run ./cmd/${PROJECT_NAME}/main.go

.PHONY: docs
docs:
	@echo ">>> Generate Swagger API Documentation..."
	swag init --generalInfo cmd/${PROJECT_NAME}/main.go

.PHONY: clean
clean:
	@echo ">>> Removing old binaries and env files..."
	@rm -rf bin/*
	@rm -rf .env
