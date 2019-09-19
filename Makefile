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
BIN_DIR_WINDOWS_X86_64=$(BIN_DIR)/windows_x86_64
BIN_DIR_LINUX_X86_64=$(BIN_DIR)/linux_x86_64
BINARY_LINUX_X86_64=$(BIN_DIR_LINUX_X86_64)/$(BINARY_NAME)
BIN_DIR_LINUX_ARM64=$(BIN_DIR)/linux_arm64
BINARY_LINUX_ARM64=$(BIN_DIR_LINUX_ARM64)/$(BINARY_NAME)
BINARY_WINDOWS_X86_64=$(BIN_DIR_WINDOWS_X86_64)/$(BINARY_NAME).exe
BINARY_LOCAL=$(BIN_DIR_LOCAL)/$(BINARY_NAME)
DIST_DIR=$(BUILD_DIR)/dist
DIST_DIR_LINUX_X86_64=$(DIST_DIR)/linux_x86_64
DIST_DIR_LINUX_ARM64=$(DIST_DIR)/linux_arm64
DIST_DIR_WINDOWS_X86_64=$(DIST_DIR)/windows_x86_64
CMD=./cmd/$(BINARY_NAME)/ifpl.go
VERSION=$(shell git describe --tags)

all: test build
build:
	mkdir -p $(BIN_DIR_LOCAL)
	$(GOBUILD) -o $(BINARY_LOCAL) -v $(CMD)
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -rf $(BIN_DIR)
	rm -rf $(DIST_DIR)

# Cross compilation
build-linux-x86-64:
	mkdir -p $(BIN_DIR_LINUX_X86_64)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX_X86_64) -v $(CMD)

dist-linux-x86-64: build-linux-x86-64
	mkdir -p $(DIST_DIR_LINUX_X86_64)
	cp README.md $(DIST_DIR_LINUX_X86_64)/
	cp LICENSE $(DIST_DIR_LINUX_X86_64)/
	cp $(BINARY_LINUX_X86_64) $(DIST_DIR_LINUX_X86_64)/
	cd $(DIST_DIR_LINUX_X86_64) && tar -zcvf ../$(BINARY_NAME)_$(VERSION)_linux_x86_64.tar.gz * && cd -

build-linux-arm64:
	mkdir -p $(BIN_DIR_LINUX_ARM64)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BINARY_LINUX_ARM64) -v $(CMD)

dist-linux-arm64: build-linux-arm64
	mkdir -p $(DIST_DIR_LINUX_ARM64)
	cp README.md $(DIST_DIR_LINUX_ARM64)/
	cp LICENSE $(DIST_DIR_LINUX_ARM64)/
	cp $(BINARY_LINUX_ARM64) $(DIST_DIR_LINUX_ARM64)/
	cd $(DIST_DIR_LINUX_ARM64) && tar -zcvf ../$(BINARY_NAME)_$(VERSION)_linux_arm64.tar.gz * && cd -

build-windows-x86-64:
	mkdir -p $(BIN_DIR_WINDOWS_X86_64)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS_X86_64) -v $(CMD)

dist-windows-x86-64: build-windows-x86-64
	mkdir -p $(DIST_DIR_WINDOWS_X86_64)
	cp README.md $(DIST_DIR_WINDOWS_X86_64)/
	cp LICENSE $(DIST_DIR_WINDOWS_X86_64)/
	cp $(BINARY_WINDOWS_X86_64) $(DIST_DIR_WINDOWS_X86_64)/
	cd $(DIST_DIR_WINDOWS_X86_64) && zip ../$(BINARY_NAME)_$(VERSION)_windows_x86_64.zip * && cd -

dist-all: dist-windows-x86-64 dist-linux-x86-64 dist-linux-arm64
