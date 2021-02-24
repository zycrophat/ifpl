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
BINARY_LOCAL=$(BIN_DIR_LOCAL)/$(BINARY_NAME)
DIST_DIR=$(BUILD_DIR)/dist
GOOSS=windows linux
GO_LINUX_ARCHS=386 amd64 arm arm64 ppc64 ppc64le mips  mipsle mips64 mips64le riscv64 s390x
GO_WINDOWS_ARCHS=386 amd64 arm
CMD=./cmd/$(BINARY_NAME)/ifpl.go
VERSION=$(shell git describe --tags)

DIST_LINUX_ALL=
DIST_WINDOWS_ALL=

define build_linux_template =
build-linux-$(1):
	mkdir -p $(BIN_DIR)/linux_$(1)
	CGO_ENABLED=0 GOOS=linux GOARCH=$(1) $(GOBUILD) -o $(BIN_DIR)/linux_$(1)/$(BINARY_NAME) -v $(CMD)
endef

define dist_linux_template =
dist-linux-$(1): build-linux-$(1)
	mkdir -p $(DIST_DIR)/linux-$(1)
	cp README.md $(DIST_DIR)/linux-$(1)/
	cp LICENSE $(DIST_DIR)/linux-$(1)/
	cp $(BIN_DIR)/linux_$(1)/$(BINARY_NAME) $(DIST_DIR)/linux-$(1)/
	cd $(DIST_DIR)/linux-$(1)/ && tar -zcvf ../$(BINARY_NAME)_$(VERSION)_linux_$(1).tar.gz * && cd -
DIST_LINUX_ALL += dist-linux-$(1)
endef

define build_windows_template =
build-windows-$(1):
	mkdir -p $(BIN_DIR)/windows_$(1)
	GOOS=windows GOARCH=$(1) $(GOBUILD) -o $(BIN_DIR)/windows_$(1)/$(BINARY_NAME).exe -v $(CMD)
endef

define dist_windows_template =
dist-windows-$(1): build-windows-$(1)
	mkdir -p $(DIST_DIR)/windows-$(1)
	cp README.md $(DIST_DIR)/windows-$(1)/
	cp LICENSE $(DIST_DIR)/windows-$(1)/
	cp $(BIN_DIR)/windows_$(1)/$(BINARY_NAME).exe $(DIST_DIR)/windows-$(1)/
	cd $(DIST_DIR)/windows-$(1)/ && zip ../$(BINARY_NAME)_$(VERSION)_windows_$(1).zip * && cd -
DIST_WINDOWS_ALL += dist-windows-$(1)
endef

$(foreach arch,$(GO_LINUX_ARCHS),$(eval $(call build_linux_template,$(arch))))
$(foreach arch,$(GO_LINUX_ARCHS),$(eval $(call dist_linux_template,$(arch))))
$(foreach arch,$(GO_WINDOWS_ARCHS),$(eval $(call build_windows_template,$(arch))))
$(foreach arch,$(GO_WINDOWS_ARCHS),$(eval $(call dist_windows_template,$(arch))))

$(eval dist-linux: $(DIST_LINUX_ALL))
$(eval dist-windows: $(DIST_WINDOWS_ALL))

dist-all: dist-linux dist-windows

.PHONY: dist-all all build test clean
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
