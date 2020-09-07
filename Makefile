.DEFAULT_GOAL := help
SHELL := /bin/bash

.PHONY: help
help: ## Print documentation
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: binary
binary: ## Build the go binary
	mkdir -p bin
	cd server-src && go build -o ../bin/ninja-proxy

.PHONY: docker
docker: ## Build the docker image
	docker build -t "ninja-proxy:latest" .
