prepare:
	@echo "Preparing required libraries..."
	@command -v pre-commit >/dev/null 2>&1 || brew install pre-commit
	@command -v jq >/dev/null 2>&1 || brew install jq
	@command -v wiremock-standalone >/dev/null 2>&1 || brew install wiremock-standalone
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
	go install golang.org/x/tools/cmd/goimports@latest

	@echo "Initializing environment for git..."
	pre-commit install

rehooks:
	@echo "Clear git hooks cache..."
	pre-commit clean

	@echo "Reinitializing environment for git..."
	pre-commit install

run:
	@echo "Running local development..."
	go run .

dev.up:
	@echo "Starting local development..."
	$(MAKE) local.env.prepare
	docker compose -f docker-compose.dev.yaml up

dev.down:
	@echo "Removing local development container..."
	docker-compose down

dev.build.push.local:
	@echo "Building container for pushing to local registry..."
	$(eval SERVICE=$(shell basename $(CURDIR)))
	docker build -f ./build/docker/Dockerfile -t $(SERVICE) .
	docker tag $(SERVICE) localhost:5001/$(SERVICE)
	docker push localhost:5001/$(SERVICE)

dev.up.alpha:
	@echo "Starting alpha development..."
	$(MAKE) alpha.env.prepare
	docker compose -f docker-compose.dev.yaml up --scale mysql=0 --scale redis=0 --scale kafka=0 --scale zookeeper=0 --scale redpanda-console=0

local.env.prepare:
	@echo "Preparing .env.generated for local envs..."
	cat .env.local > .env.generated
	cat .env.local.common >> .env.generated

alpha.env.prepare:
	@echo "Preparing .env.generated for alpha envs..."
	sh ./scripts/get-k8s-alpha-env.sh > .env.generated
	cat .env.local >> .env.generated

wiremock:
	wiremock --port 1324 --verbose --local-response-templating --root-dir tools/wiremock

swagger:
	@echo "Generating swagger docs..."
	swag init

.PHONY: prepare rehooks run dev.up dev.down dev.build.push.local dev.up.alpha local.env.prepare alpha.env.prepare wiremock swagger
