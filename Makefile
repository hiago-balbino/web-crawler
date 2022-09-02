.PHONY: help vet tests mongo-up mongo-down docker-ps

## help: show this help.
help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

## vet: run the command vet from Go
vet:
	go vet ./...

## tests: run all unit tests
tests:
	go test -race ./...

## docker-ps: list all containers running
docker-ps:
	docker-compose -f docker/docker-compose.yml ps

## mongo-up: start mongo container
mongo-up:
	docker-compose -f docker/docker-compose.yml up -d

## mongo-down: stop mongo container
mongo-down:
	docker-compose -f docker/docker-compose.yml down