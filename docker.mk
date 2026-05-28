.PHONY: docker-build docker-up docker-down docker-logs docker-clean

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

docker-up: ## Start all services
	@echo "Starting all services..."
	@docker compose up --build -d

docker-down: ## Stop all services
	@echo "Stopping all services..."
	@docker compose down

docker-logs: ## View logs
	@echo "Viewing logs..."
	@docker compose logs -f

docker-clean: ## Remove containers and volumes
	@echo "Removing containers and volumes..."
	@docker compose down -v --remove-orphans
