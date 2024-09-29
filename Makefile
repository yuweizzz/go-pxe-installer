GOFILE := $(shell find . -name "*.go" | xargs)

default: build

.PHONY: lint
lint:
	gofmt -w $(GOFILE)

.PHONY: build
build:
	go build
