.PHONY: fmt generate test lint run help grow dockerize-dist mock test-server
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
GOFMT := "goimports"

grow: ## Change project seed name
	find . '(' -name '*.go' -or -name 'go.mod' -or -name '*.md' ')' -exec sed -i '' -e 's/seller/$(app_name)/g' '{}' \;

fmt: ## Run gofmt for all .go files
	@$(GOFMT) -w $(GOFMT_FILES)

generate: mock ## Generate stuffs


test: ## Run go test for whole project
	@go test -v ./... || true

lint: ## Run linter
	@golangci-lint run ./...

lint-fix: ## Run linter and autofix
	@golangci-lint run --fix ./...

build-dist: ## Run builder
	GOOS=linux GOARCH=${ARCH} CGO_ENABLED=0 go build  -ldflags "-s -w" -o  ./dist/$(app_name) ./cmd/$(app_name)/*.go

dockerize-dist:  ## Dockerize dist
	@echo "Building docker image for ${DOCKER_REGISTRY}/${image_name}:${tag}"
	docker build  . -t "${DOCKER_REGISTRY}/${image_name}:${tag}" --platform linux/${ARCH}
	docker push "${DOCKER_REGISTRY}/${image_name}:${tag}"

run: ## Run server, hot reload if added args --hot
	clear
	@echo "Running applications with parameter: hot=${hot} extra=${extra}"
	@if [ -z $(hot) ]; then go run cmd/${app_name}/*.go ${extra}; else air --build.cmd "go build -o tmp/${app_name} cmd/${app_name}/*.go" --build.bin "./tmp/${app_name}" --build.args_bin "${extra}"; fi

migrate:
	$(call process_sql_files,extension)
	atlas migrate hash
	atlas migrate diff --env gorm
	$(call process_sql_files,trigger)

# Shared function to process SQL files
define process_sql_files
	@if [ -d "sql/$(1)s" ]; then \
		files="$$(find sql/$(1)s -name '*.sql' | (sort -V 2>/dev/null || sort))"; \
		if [ -n "$$files" ]; then \
			echo "$$files" | xargs cat > tmp.sql; \
			latest="$$(find migrations -name "*_$(1).sql" 2>/dev/null | (sort -Vr 2>/dev/null || sort -r) | head -n 1)"; \
			if [ -f "$$latest" ] && diff -q "$$latest" tmp.sql; then \
				rm tmp.sql; \
			else \
				sleep 1; \
				diff "$$latest" tmp.sql | head -n 10; \
				echo "New $(1) file created"; \
				mv tmp.sql migrations/$$(TZ=UTC date +"%Y%m%d%H%M%S")_$(1).sql; \
				atlas migrate hash; \
			fi; \
		else \
			echo "No '.sql' files found in 'sql/$(1)s' directory"; \
		fi; \
	else \
		echo "No 'sql/$(1)s' directory"; \
	fi
endef

mock:
	clear
	@mockery

test-server:
	clear
	@goconvey

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | LC_ALL=en_US,UTF-8 sort -s | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
