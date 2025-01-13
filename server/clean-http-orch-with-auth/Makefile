include .env.local
export

## prepare: prepare required libraries for development
.PHONY: prepare
prepare:
	@echo "Preparing required libraries..."
	@command -v pre-commit >/dev/null 2>&1 || brew install pre-commit
	@command -v jq >/dev/null 2>&1 || brew install jq
	@command -v teleport >/dev/null 2>&1 || { go install github.com/project-inari/teleport@latest; }
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
	go install golang.org/x/tools/cmd/goimports@latest

	@echo "Initializing environment for git..."
	pre-commit install

## rehooks: reinitialize git hooks
.PHONY: rehooks
rehooks:
	@echo "Clear git hooks cache..."
	pre-commit clean

	@echo "Reinitializing environment for git..."
	pre-commit install

## upgrade-libs: run upgrade all libraries
.PHONY: upgrade-libs
upgrade-libs:
	go get -t -u ./...
	go mod tidy

## golangci-lint: run golangci-lint
.PHONY: golangci-lint
golangci-lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run ./...

## run: run the application
.PHONY: run
run:
	@echo "Running local development..."
	go run .

## dev.up: start local development in docker
.PHONY: dev.up
dev.up:
	@echo "Starting local development..."
	$(MAKE) local.env.prepare
	docker compose -f docker-compose.dev.yaml up --scale wiremock=0

## dev.down: stop local development in docker
.PHONY: dev.down
dev.down:
	@echo "Removing local development container..."
	docker-compose down

## dev.build.push.local: build and push container to registry
.PHONY: dev.build.push
dev.build.push:
	@echo "Building container for pushing to registry..."
	$(eval SERVICE=$(shell basename $(CURDIR)))
	docker build -f ./build/docker/Dockerfile -t $(SERVICE) .
	docker tag $(SERVICE) asia-southeast1-docker.pkg.dev/inari-non-prod/inari-registry/$(SERVICE)
	docker push asia-southeast1-docker.pkg.dev/inari-non-prod/inari-registry/$(SERVICE)

## dev.up.alpha: start local development in docker with alpha environment
.PHONY: dev.up.alpha
dev.up.alpha:
	@echo "Starting alpha development..."
	$(MAKE) alpha.env.prepare
	sh -c '$(MAKE) portforward' 2>&1 &
	docker compose -f docker-compose.dev.yaml up --scale mysql=0 --scale redis=0 --scale wiremock=0

## local.env.prepare: prepare .env.generated for local envs
.PHONY: local.env.prepare
local.env.prepare:
	@echo "Preparing .env.generated for local envs..."
	cat .env.local > .env.generated
	cat .env.local.common >> .env.generated

## alpha.env.prepare: prepare .env.generated for alpha envs
.PHONY: alpha.env.prepare
alpha.env.prepare:
	@echo "Preparing .env.generated for alpha envs..."
	sh ./scripts/get-k8s-alpha-env.sh > .env.generated
	cat .env.dev.alpha >> .env.generated

## wiremock: start local wiremock server
.PHONY: wiremock
wiremock:
	@command -v wiremock-standalone >/dev/null 2>&1 || brew install wiremock-standalone
	wiremock --port 1324 --verbose --local-response-templating --root-dir tools/wiremock

## swagger: generate swagger docs
.PHONY: swagger
swagger:
	@echo "Generating swagger docs..."
	swag init

## portforward: portforward to k8s cluster from tpconfig
.PHONY: portforward
portforward:
	teleport -s .tpconfig

## test: run all tests
.PHONY: test
test: tidy
	go test ./...

## test-cover: run all tests and display coverage
.PHONY: test-cover
test-cover:
	go test -v -race -buildvcs -coverprofile=./coverage.out ./...
	go tool cover -html=./coverage.out
