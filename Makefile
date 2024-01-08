##@ General
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

##@ Build
build: fmt vet  ## Build service binary.
	go build -o bin/lambda main.go

run: vet ## Run service from your laptop.
	go run ./main.go

##@ Test
lint: ## Run Go linter
	golangci-lint run ./...

test: ## Run Go tests
	go test ./...