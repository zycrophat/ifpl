SHELL:=/bin/bash

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=ifpl
BUILD_DIR=./build
BIN_DIR=$(BUILD_DIR)/bin
BIN_DIR_LOCAL=$(BIN_DIR)/local
BIN_DIR_WINDOWS=$(BIN_DIR)/windows
BIN_DIR_LINUX=$(BIN_DIR)/linux
BINARY_LINUX=$(BIN_DIR_LINUX)/$(BINARY_NAME)
BINARY_WINDOWS=$(BIN_DIR_WINDOWS)/$(BINARY_NAME).exe
BINARY_LOCAL=$(BIN_DIR_LOCAL)/$(BINARY_NAME)
DIST_DIR=$(BUILD_DIR)/dist
DIST_DIR_LINUX=$(DIST_DIR)/linux
DIST_DIR_WINDOWS=$(DIST_DIR)/windows
CMD=./cmd/$(BINARY_NAME)/ifpl.go
VERSION=$(shell git describe --tags)

all: test build
build:
	mkdir -p $(BUILD_DIR_LOCAL)
	$(GOBUILD) -o $(BINARY_LOCAL) -v $(CMD)
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -rf $(BIN_DIR)
	rm -rf $(DIST_DIR)

# Cross compilation
build-linux:
	mkdir -p $(BIN_DIR_LINUX)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX) -v $(CMD)

dist-linux: build-linux
	mkdir -p $(DIST_DIR_LINUX)
	cp README.md $(DIST_DIR_LINUX)/
	cp LICENSE $(DIST_DIR_LINUX)/
	cp $(BINARY_LINUX) $(DIST_DIR_LINUX)/
	cd $(DIST_DIR_LINUX) && tar -zcvf ../$(BINARY_NAME)_$(VERSION)_linux_x86_64.tar.gz * && cd -

build-windows:
	mkdir -p $(BIN_DIR_WINDOWS)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS) -v $(CMD)

dist-windows: build-windows
	mkdir -p $(DIST_DIR_WINDOWS)
	cp README.md $(DIST_DIR_WINDOWS)/
	cp LICENSE $(DIST_DIR_WINDOWS)/
	cp $(BINARY_WINDOWS) $(DIST_DIR_WINDOWS)/
	cd $(DIST_DIR_WINDOWS) && zip ../$(BINARY_NAME)_$(VERSION)_windows_x86_64.zip * && cd -

dist-all: dist-windows dist-linux
