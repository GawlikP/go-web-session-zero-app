.PHONY: lint

lint:
	go vet ./...
	go fmt ./...
generate:
	templ generate ./...

build:
	go build -o ssr ./cmd/ssr

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development commands
dev:
	podman-compose -f podman-compose-development.yml up

dev-build:
	podman-compose -f podman-compose-development.yml up --build

dev-down:
	podman-compose -f podman-compose-development.yml down

dev-logs:
	podman-compose -f podman-compose-development.yml logs -f

# Production commands
prod:
	podman-compose up -d

prod-build:
	podman-compose up -d --build

prod-down:
	podman-compose down

prod-logs:
	podman-compose logs -f

# Utility commands
clean:
	podman system prune -a -f
	podman volume prune -f:w!
