PROJECT_NAME := "dns"
PKG := "github.com/cmol/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)
EXAMPLES := $(shell ls examples)

.PHONY: all dep build clean test coverage coverhtml lint examples

all: lint build

lint: ## Lint the files
	@golangci-lint run

test: ## Run unittests
	@go test -short ${PKG_LIST}

race: dep ## Run data race detector
	@go test -race -short ${PKG_LIST}

msan: dep ## Run memory sanitizer
	@env CC=clang env CXX=clang++ go test -msan -short ${PKG_LIST}

coverage: ## Generate global code coverage report
	@go test -cover ${PKG_LIST}

coverhtml: ## Generate global code coverage report in HTML
	@go test -coverprofile=coverage.out && go tool cover -html=coverage.out

dep: ## Get the dependencies
	@go get -v -d ./...

build: dep ## Build the binary file
	@go build -v $(PKG)

examples: dep ## Build examples
	@for example in $(EXAMPLES); do go build ./examples/$$example; done

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME) $(EXAMPLES)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
