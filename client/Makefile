CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s -w' -o your_binary_name
# Variables
APP_NAME = client
GO_CMD = go
GO_BUILD = $(GO_CMD) build
GO_TEST = $(GO_CMD) test

# Default target: build
all: build

# Find all Go files in the project
GO_FILES := $(shell find . -name '*.go')

# Build target: compile the Go source files into the final executable
build: $(APP_NAME)

clean:
	rm bin/${APP_NAME}

# Compile the Go source files
$(APP_NAME): $(GO_FILES)
# rm bin/${APP_NAME}
	CGO_ENABLED=0
	$(GO_BUILD) -tags netgo -ldflags '-extldflags "-static"'  -o bin/$(APP_NAME)

# Test target: run the tests for the project
test:
	$(GO_TEST) -v ./... -cover

# Declare phony targets
.PHONY: all build test