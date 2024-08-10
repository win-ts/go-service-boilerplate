prepare:
	@echo "Preparing required libraries..."
	brew install pre-commit
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
	go install golang.org/x/tools/cmd/goimports@latest

	@echo "Initializing environment for git..."
	pre-commit install -t pre-commit
	pre-commit install -t pre-push
	pre-commit install -t commit-msg

rehooks:
	@echo "Clear git hooks cache..."
	pre-commit clean

	@echo "Reinitializing environment for git..."
	pre-commit install -t pre-commit
	pre-commit install -t pre-push
	pre-commit install -t commit-msg

.PHONY: prepare rehooks
