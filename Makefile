include .env
export

.PHONY: run run-in build clean db app down doc test lint

run: ## Start the application in detached mode
	docker-compose up -d

run-debug:
	docker-compose -f local-docker-compose.yaml up

run-in: ## Start the application in interactive mode
	docker-compose up

build: ## Build Docker images
	docker-compose build

clean: ## Remove unused Docker images
	docker system prune -af
	docker volume prune -af
	docker system df
	docker rmi -f $$(docker images -q) || true

db: ## Connect to the database
	docker exec -it 2023_2_hamster-db-1 psql -U $(DB_USER) -d $(DB_NAME)

app: ## Connect to the application container
	docker exec -it 2023_2_hamster-server-1 ./app

down: ## Stop and remove containers, networks, images, and volumes
	docker-compose down

doc: ## Generate API documentation using swag
	swag init -g cmd/api/main.go

test: ## Run tests
	# Add your test command here

lint: ## Run linters
	golangci-lint run

# New make

# b:
# 	go build -o app ./cmd/api/main.go

# r: lint b	
# 	./app