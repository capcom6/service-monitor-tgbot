.PHONY: swagger

swagger: ## Generate Swagger documentation
	@if ! swag fmt -g ./main.go; then \
		echo "Error: Failed to format API docs"; \
		exit 1; \
	fi
	@if ! swag init --parseDependency --outputTypes go -g ./main.go -o ./internal/server/docs; then \
		echo "Error: Failed to generate API docs"; \
		exit 1; \
	fi
