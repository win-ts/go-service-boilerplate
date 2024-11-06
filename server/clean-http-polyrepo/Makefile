dev.up:
	@echo "Starting local development..."
	docker compose -f docker-compose.dev.yaml up

dev.down:
	@echo "Removing local development container..."
	docker-compose down

dev.build.push:
	@echo "Building container for pushing to registry..."
	$(eval SERVICE=$(shell basename $(CURDIR)))
	docker build -f ./build/docker/Dockerfile -t $(SERVICE) --build-arg SERVICE=$(SERVICE) .
	docker tag $(SERVICE) localhost:5001/$(SERVICE)
	docker push localhost:5001/$(SERVICE)

swagger:
	@echo "Generating swagger docs..."
	swag init

.PHONY: dev.up dev.down dev.build.push swagger
