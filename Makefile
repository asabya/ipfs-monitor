GOCMD := go
GOBUILD := $(GOCMD) build
BUILD_DIR := $(PWD)
BIN_PATH := $(BUILD_DIR)/cmd
BINARY_NAME := mipfs
BUILD_SUPPORT := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))/.build-support
VERSION := "v0.1.0-dev"

build:
	@-echo "Building binary..."
	@-$(GOBUILD) -o $(BIN_PATH)/$(BINARY_NAME)

install:
	@-echo "Installing binary..."
	@-$(GOBUILD) -o $(GOBIN)/$(BINARY_NAME)

all:
	@-echo "Building binaries for all platforms"
	@. $(BUILD_SUPPORT) ; BIN_PATH=$(BIN_PATH) BINARY_NAME=$(BINARY_NAME) makeall $(VERSION)
