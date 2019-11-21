# list of pkgs for the project without vendor
PKGS=$(shell go list ./... | grep -v /vendor/)

APP_NAME=playbtg

help: ## Print help commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean: ## Remove working files and binary file
	@go clean
	@rm -Rf dist

build: ## Build the binary file
	@CGO_ENABLED=0 GOGC=off GOOS=linux GOARCH=amd64 go build -v -a -installsuffix nocgo -o dist/$(APP_NAME) main.go

docker-build: ## Build docker image with latest tag
	@docker build --tag=$(APP_NAME):latest .

format: ## Format go code
	@go fmt $(PKGS)
	@goimports -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

test: ## Run unit tests
	@go test -v $(PKGS)

start: ## Start the app
	@go run main.go

run-bin: build ## Start the app from binary file
	@dist/$(APP_NAME)