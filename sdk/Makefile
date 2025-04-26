.PHONY: fmt generate test lint run help grow dockerize-dist mock test-server
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
GOFMT := "goimports"

fmt: ## Run gofmt for all .go files
	@$(GOFMT) -w $(GOFMT_FILES)

test: ## Run go test for whole project
	@go test -v ./...

lint: ## Run linter
	@golangci-lint run ./...

lint-fix: ## Run linter with auto-fix
	@golangci-lint run --fix ./...

generate: mock ## Generate stuffs

mock:
	clear
	@mockery

tidy:
	@go mod tidy

test-server:
	clear
	@goconvey

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
