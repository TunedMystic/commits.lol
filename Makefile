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

clean:  ## Clean workspace
	@rm -f ${APP}
	@rm -rf coverage.out
	@go clean -testcache

dev:  ## Run the program in dev mode.
	@go run main.go

install:  ## Install project dependencies
	@go mod download

test: clean  ## Run tests
	@eval $$(egrep -v '^#' .env.test | xargs) go test ./... -covermode=atomic -coverprofile coverage.out
	@go tool cover -func coverage.out
	@eval $$(egrep -v '^#' .env.test | xargs) bash scripts/coverage-threshold.sh

cover:  ## View HTML coverage reports
	@go tool cover -html coverage.out

watch:  ## Watch for file changes and run the server.
	@npm run build-styles
	@bash -c "find . -type f \( -name '*.go' -o -name '*.html' \) | grep -v 'misc' | entr -r $(MAKE) dev"

watchtests:  ## Watch for file changes and run tests.
	@bash -c "find . -name '*.go' | grep -v 'misc' | entr -r $(MAKE) test"

.PHONY: help build clean dev install test watch watchtests
