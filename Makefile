BASE_DIR := $(shell pwd)
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

.PHONY: ipxe
ipxe:
	@$(BASE_DIR)/utils/ipxe.sh

.PHONY: images
images:
	@$(BASE_DIR)/utils/images.sh

.PHONY: build
build:
	CGO_ENABLED=0 go build -o $(BIN_DIR)/$(BIN_NAME)

.PHONY: buildi
buildi:
	mv help/images tftpboot
	mv help/debian12-preseed.txt tftpboot
	CGO_ENABLED=0 go build -o $(BIN_DIR)/$(BIN_NAME)
	mv tftpboot/images $(BASE_DIR)/help
	mv tftpboot/debian12-preseed.txt $(BASE_DIR)/help

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
