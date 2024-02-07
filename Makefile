.PHONY: all help setup vet lint deadcode vulncheck tests integration-tests all-tests cover sonarqube-up sonarqube-down sonarqube-analysis fmt compose-ps compose-up compose-down build build-run-api clean doc

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

## lint: run all linters configured
lint:
	golangci-lint run ./...	

## deadcode: run command to look for deadcode
deadcode:
	deadcode .

## vulncheck: run all vulnerability checks
vulncheck:
	govulncheck ./...

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

## sonarqube-up: start sonarqube container
sonarqube-up:
	docker run -d --name sonarqube -p ${SONAR_PORT}:${SONAR_PORT} sonarqube

## sonarqube-down: stop sonarqube container
sonarqube-down:
	docker rm sonarqube -f

## sonarqube-analysis: run sonar scanner
sonarqube-analysis: all-tests
	${SONAR_BINARY} -Dsonar.host.url=${SONAR_HOST} -Dsonar.login=${SONAR_LOGIN} -Dsonar.password=${SONAR_PASSWORD}

## fmt: run go formatter recursively on all files
fmt:
	gofmt -s -w .

## compose-ps: list all containers running
compose-ps:
	docker-compose -f build/docker-compose.yml ps

## compose-up: start API and dependencies
compose-up:
	docker-compose -f build/docker-compose.yml up -d

## compose-down: stop API and dependencies
compose-down:
	docker-compose -f build/docker-compose.yml down

## build: create an executable of the application
build:
	go build -o ${APP_NAME} .

## build-run-api: build project and run the API using the built binary
build-run-api: build
	./${APP_NAME} api

## clean: run the go clean command and removes the application binary
clean:
	go clean
	rm ${APP_NAME}

## doc: run the project documentation using HTTP
doc:
	godoc -http=:6060

all: help
