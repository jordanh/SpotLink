# Define variables
BINARY_NAME=spotlink
BUILD_DIR=./bin
ZIP_FILE=function.zip

# Default target executed when no arguments are given to make.
default: build

# Build binary for AWS Lambda deployment
build:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	chmod +x $(BUILD_DIR)/$(BINARY_NAME)

# Package binary into a zip file for AWS Lambda deployment
package: build
	zip -j $(BUILD_DIR)/$(ZIP_FILE) $(BUILD_DIR)/$(BINARY_NAME)

# Clean up the binary and zip file
clean:
	rm -rf $(BUILD_DIR)/*

# Deploy to AWS Lambda (replace <function-name> with your actual function name)
deploy: package
	aws lambda update-function-code --function-name GoSpotLink --zip-file fileb://$(BUILD_DIR)/$(ZIP_FILE)

# Phony targets
.PHONY: default build package clean deploy
