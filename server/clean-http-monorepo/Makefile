dev.up:
	@echo "Starting local development..."
	docker compose -f docker-compose.dev.yaml up

dev.down:
	@echo "Removing local development container..."
	docker-compose down

dev.build.push:
	@echo "Building container for pushing to registry..."
	$(eval DOMAIN=$(shell basename $(shell dirname $(CURDIR))))
	$(eval SERVICE=$(shell basename $(CURDIR)))
	docker build -f ../../build/docker/Dockerfile -t $(SERVICE) --build-arg SERVICE=$(DOMAIN)/$(SERVICE) ../../
	docker tag $(SERVICE) asia-southeast1-docker.pkg.dev/inari-non-prod/inari-registry/$(SERVICE)
	docker push asia-southeast1-docker.pkg.dev/inari-non-prod/inari-registry/$(SERVICE)

k8s.dev:
	@echo "Deploying to local k8s in dev namespace..."
	kubectl apply -k ./k8s/overlays/dev

swagger:
	@echo "Generating swagger docs..."
	swag init

.PHONY: dev.up dev.down dev.build.push k8s.dev swagger
