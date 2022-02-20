GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

NAME=nazo-api
CMD_DIR=cmd/nazo
BUILD_DIR=build
BUILD_DOC_DIR=$(BUILD_DIR)/doc

define DOCKERFILE_REDOC
FROM node:latest
RUN npm install -g redoc-cli && npm install -g @redocly/openapi-cli
endef
export DOCKERFILE_REDOC

.PHONY: build build-all clean build-all-on-docker

build:
	@GO111MODULE=on go build -ldflags '-s -w' -o $(BUILD_DIR)/$(GOOS)_$(GOARCH)/$(NAME) ./$(CMD_DIR)

build-all: build $(BUILD_DOC_DIR)/index.html

clean:
	@rm -rf $(BUILD_DIR)

build-all-on-docker:
	@echo "$$DOCKERFILE_REDOC" | docker build --quiet -t redoc -
	@docker run -v $(PWD):/src -w /src --rm -i redoc /usr/bin/make $(BUILD_DOC_DIR)/index.html
	docker run -v $(PWD):/src -w /src --rm -i golang /usr/bin/make build


$(BUILD_DOC_DIR)/index.html: $(BUILD_DOC_DIR)/openapi.yaml
	@mkdir -p $(BUILD_DOC_DIR)
	@redoc-cli bundle $(BUILD_DOC_DIR)/openapi.yaml -o $(BUILD_DOC_DIR)/index.html

$(BUILD_DOC_DIR)/openapi.yaml:
	@mkdir -p $(BUILD_DOC_DIR)
	@openapi bundle api/openapi/openapi.yaml >$(BUILD_DOC_DIR)/openapi.yaml
