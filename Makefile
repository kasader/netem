.DEFAULT_GOAL := help

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: lint
lint: ## Runs golangci-lint for static code analysis.
	golangci-lint run ./..

.PHONY: docs
docs: ## Start the docs webserver using golang.org/x/pkgsite.
	go doc -http