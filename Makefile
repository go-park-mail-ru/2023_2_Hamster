include .env
export

.PHONY: run run-in build clean db app down doc test lint

run: ## Start the application in detached mode
	docker-compose up -d

run-in: ## Start the application in interactive mode
	docker-compose up

build: ## Build Docker images
	docker-compose build

clean: ## Remove unused Docker images
	docker rmi -f $$(docker images -q)

db: ## Connect to the database
	docker exec -it $(PROJECT_NAME)_db-1 psql -U $(DB_USER) -d $(DB_NAME)

app: ## Connect to the application container
	docker exec -it $(PROJECT_NAME)_server-1 ./app

down: ## Stop and remove containers, networks, images, and volumes
	docker-compose down

doc: ## Generate API documentation using swag
	swag init -g cmd/api/main.go

test: ## Run tests
	# Add your test command here

lint: ## Run linters
	golangci-lint run
