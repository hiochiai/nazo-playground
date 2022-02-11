GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

NAME=nazo-api
CMD_DIR=cmd/nazo
BUILD_DIR=build

.PHONY: build
build:
	@GO111MODULE=on go build -ldflags '-s -w' -o $(BUILD_DIR)/$(GOOS)_$(GOARCH)/$(NAME) $(PWD)/$(CMD_DIR)

clean:
	@rm -rf $(BUILD_DIR)