
CLI_SOURCE_FILES?=./cmd/plugin
CLI_BINARY_NAME=binary
CLI_DESTINATION=./bin/$(CLI_BINARY_NAME)
CLI_DESTINATION_LOCAL_OSX="${HOME}/Library/Application Support/atlascli/plugins/mmarcon@cli-plugin-tools/$(CLI_BINARY_NAME)"

.PHONY: build
build: ## Generate the binary in ./bin
	@echo "==> Building $(CLI_BINARY_NAME) binary"
	go build -o $(CLI_DESTINATION) $(CLI_SOURCE_FILES)

.PHONY: build-local-osx
build-local-osx:
	@echo "==> Building $(CLI_DESTINATION_LOCAL_OSX)"
	go build -o $(CLI_DESTINATION_LOCAL_OSX) $(CLI_SOURCE_FILES)

.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'