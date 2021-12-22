.PHONY: help vet tests

## help:					show this help.
help:
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

## vet:					run the command vet from Go
vet:
	go vet ./...

## tests:					run all unit tests
tests:
	go test -race ./...