APP=commits.lol

help:  ## This help
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[1m%-15s\033[0m %s\n", $$1, $$2}'

build: clean  ## Build the binary
	@npm run build-styles-prod
	@go build -ldflags="-s -w"

build-docker: clean  ## Build docker image
	@npm run build-styles-prod
	@docker build -t commits.lol .

install:  ## Install project dependencies
	@go mod download

clean:  ## Clean workspace
	@rm -f ${APP}
	@rm -f coverage.out
	@go clean -testcache

dev:  ## Run the program
	@eval $$(egrep -v '^#' .env | xargs) go run main.go server

dev-w:  ## Run the program and watch for file changes
	@npm run build-styles
	@bash -c "find . -type f \( -name '*.go' -o -name '*.html' \) | grep -v 'misc' | entr -r $(MAKE) dev"

dev-docker: ## Run the docker container.
	@docker run -p 8000:8000 --env-file .env -v $$(pwd)/commits.lol.sqlite:/usr/src/commits.lol.sqlite commits.lol

test: clean  ## Run tests
	@eval $$(egrep -v '^#' .env.test | xargs) go test ./... -covermode=atomic -coverprofile coverage.out
	@go tool cover -func coverage.out
	@eval $$(egrep -v '^#' .env.test | xargs) bash scripts/coverage-threshold.sh

test-w:  ## Run tests and watch for file changes
	@bash -c "find . -name '*.go' | grep -v 'misc' | entr -r $(MAKE) test"

cover:  ## View HTML coverage reports
	@go tool cover -html coverage.out

.PHONY: help build build-docker install clean dev dev-w dev-docker test test-w cover
