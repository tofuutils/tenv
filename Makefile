#
# Copyright 2024 tofuutils authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

.PHONY: build test clean

##@ General
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

fmt: ## Run go fmt against code.
	go fmt ./...

get: ## Install dependencies.
	go get ./...

vet: ## Run go vet against code.
	go vet ./...

##@ Build
build: get fmt ## Build service binary.
	mkdir ./build || echo
	go build -o ./build/tenv ./cmd/tenv
	go build -o ./build/tofu ./cmd/tofu
	go build -o ./build/terraform ./cmd/terraform
	go build -o ./build/terragrunt ./cmd/terragrunt
	go build -o ./build/atmos ./cmd/atmos

##@ Run
run: build ## Run service from your laptop.
	./build/tenv

##@ Lint
lint: ## Run Go linter
	golangci-lint run ./...

##@ Test
test: ## Run Go tests
	go test ./...

##@ Clean
clean:
	rm -f ./build