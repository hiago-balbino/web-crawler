.PHONY: all help setup vet tests integration-tests all-tests cover lint fmt compose-ps compose-up compose-down build build-run-api clean

APP_NAME=crawler

## help: show this help.
help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

## setup: run the command mod download and tidy from Go
setup:
	GO111MODULE=on go mod download
	go mod tidy
	go mod verify

## vet: run the command vet from Go
vet:
	go vet ./...

## tests: run all unit tests
tests:
	go test -race -coverprofile coverage.out ./... -short=true -count=1

## integration-tests: run all integration tests
integration-tests:
	go test -race ./... -run Integration -count=1

## all-tests: run all unit and integration tests
all-tests:
	go test -race -coverprofile coverage.out ./... -count=1

## cover: run the command tool cover to open coverage file as HTML
cover: all-tests
	go tool cover -html coverage.out

## lint: run all linters configured
lint:
	golangci-lint run ./...	

## fmt: run go formatter recursively on all files
fmt:
	gofmt -s -w .

## compose-ps: list all containers running
compose-ps:
	docker-compose -f docker-compose.yml ps

## compose-up: start API and dependencies
compose-up:
	docker-compose -f docker-compose.yml up -d

## compose-down: stop API and dependencies
compose-down:
	docker-compose -f docker-compose.yml down

## build: create an executable of the application
build:
	go build -o ${APP_NAME} .

## build-run-api: build project and run the API using the built binary
build-run-api: build
	./${APP_NAME} api

## clean: runs the go clean command and removes the application binary
clean:
	go clean
	rm ${APP_NAME}

all: help
