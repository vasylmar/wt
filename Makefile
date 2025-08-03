APP_NAME := wt
BUILD_DIR := build
BIN_DIR := $(HOME)/go/bin
SRC := .

.PHONY: all build clean install uninstall

all: build

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC)

install: build
	@mkdir -p $(BIN_DIR)
	@cp $(BUILD_DIR)/$(APP_NAME) $(BIN_DIR)/
	@echo "Installed $(APP_NAME) to $(BIN_DIR)"

uninstall:
	@rm -f $(BIN_DIR)/$(APP_NAME)
	@echo "Uninstalled $(APP_NAME) from $(BIN_DIR)"

clean:
	@rm -rf $(BUILD_DIR)
