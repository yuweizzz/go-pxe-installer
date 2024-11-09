BIN_DIR := bin
BIN_NAME := go-pxe-installer
GO_FILES := $(shell find . -name "*.go" | xargs)

default: build

.PHONY: lint
lint:
	gofmt -w $(GO_FILES)

.PHONY: build
build:
	go build -o $(BIN_DIR)/$(BIN_NAME)

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
