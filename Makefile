SHELL := /bin/bash

BINARY_NAME := gopud
MAIN := ./
BUILD_DIR := build

GOOS ?= linux android darwin windows
GOARCH ?= amd64 arm64

GARBLE_FLAGS := -literals -tiny -seed=random
LDFLAGS := -s -w
CGO_ENABLED ?= 0

.PHONY: all clean build linux windows darwin android native

all: clean build

build:
	@mkdir -p $(BUILD_DIR)
	@for os in $(GOOS); do \
		for arch in $(GOARCH); do \
			ext=""; \
			if [ $$os = "windows" ]; then ext=".exe"; fi; \
			printf "Building %s-$$os-$$arch$$ext ...\n" "$(BINARY_NAME)"; \
			CGO_ENABLED=$(CGO_ENABLED) GOOS=$$os GOARCH=$$arch sh -c '\
				if command -v garble >/dev/null 2>&1; then \
					garble $(GARBLE_FLAGS) build -v -trimpath -ldflags "$(LDFLAGS)" -o "$(BUILD_DIR)/$(BINARY_NAME)-$$os-$$arch$$ext" $(MAIN); \
				else \
					go build -v -trimpath -ldflags "$(LDFLAGS)" -o "$(BUILD_DIR)/$(BINARY_NAME)-$$os-$$arch$$ext" $(MAIN); \
				fi' || exit $$?; \
		done; \
	done

linux:
	$(MAKE) GOOS=linux GOARCH=amd64 build

linux-arm64:
	$(MAKE) GOOS=linux GOARCH=arm64 build

android:
	$(MAKE) GOOS=android GOARCH=amd64 build

android-arm64:
	$(MAKE) GOOS=android GOARCH=arm64 build

darwin:
	$(MAKE) GOOS=darwin GOARCH=amd64 build

darwin-arm64:
	$(MAKE) GOOS=darwin GOARCH=arm64 build

windows:
	$(MAKE) GOOS=windows GOARCH=amd64 build

windows-arm64:
	$(MAKE) GOOS=windows GOARCH=arm64 build

native:
	@mkdir -p $(BUILD_DIR)
	@printf "Building native binary (CGO_ENABLED=1)...\n"
	@CGO_ENABLED=1 sh -c '\
		if command -v garble >/dev/null 2>&1; then \
			garble $(GARBLE_FLAGS) build -v -trimpath -ldflags "$(LDFLAGS)" -o "$(BUILD_DIR)/$(BINARY_NAME)-native" $(MAIN); \
		else \
			go build -v -trimpath -ldflags "$(LDFLAGS)" -o "$(BUILD_DIR)/$(BINARY_NAME)-native" $(MAIN); \
		fi'

clean:
	rm -rf $(BUILD_DIR)
