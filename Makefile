SHELL := bash

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: tidy
lint:
	go vet ./...
	go fmt ./...

.PHONY: unit-test
unit-test:
	go test ./...

.PHONY: build
build:
	../scripts/build.sh
