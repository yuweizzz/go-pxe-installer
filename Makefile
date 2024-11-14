BIN_DIR := bin
BIN_NAME := go-pxe-installer
GO_FILES := $(shell find . -name "*.go" | xargs)

ifndef WITH_IMAGES
	WITH_IMAGES = 0
endif

default: build

.PHONY: lint
lint:
	gofmt -w $(GO_FILES)

.PHONY: build
build:
	if [ $(WITH_IMAGES) = 1 ]; then \
		mv help/images tftpboot;    \
	fi
	CGO_ENABLED=0 go build -o $(BIN_DIR)/$(BIN_NAME)
	if [ $(WITH_IMAGES) = 1 ]; then \
		mv tftpboot/images help;    \
	fi

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
