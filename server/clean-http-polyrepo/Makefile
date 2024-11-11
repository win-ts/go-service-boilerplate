prepare:
	@echo "Preparing required libraries..."
	brew install pre-commit
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

swagger:
	@echo "Generating swagger docs..."
	swag init

.PHONY: prepare rehooks run dev.up dev.down dev.build.push.local swagger
