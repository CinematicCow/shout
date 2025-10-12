BINARY_NAME=shout
BUILD_DIR=.

.PHONY: all build install clean

all: build

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/main.go

install:
	./scripts/symlink.sh

clean:
	rm -f $(BUILD_DIR)/$(BINARY_NAME)
