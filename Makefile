include .env
export

.PHONY: all run build clean db app down doc test lint cover

all: run

run: lint build up
	./server

build: ## Build Docker images && api
	docker-compose -f local-docker-compose.yaml build ; \
	go build -o server ./cmd/api/main.go

up: 
	docker-compose -f local-docker-compose.yaml up -d ; \
	sleep 5;

down: ## Stop and remove containers, networks, images, and volumes
	docker-compose -f local-docker-compose.yaml down

clean: down ## Remove unused Docker images
	docker system prune -af 
	docker volume prune -af
	docker system df
	docker rmi -f $$(docker images -q) || true
	rm -rf ./server

db: ## Connect to the database
	docker exec -it hammy-db psql -U $(DB_USER) -d $(DB_NAME)

cover:
	sh scripts/coverage_test.sh

lint: ## Run linters
	golangci-lint run

test: ## Run tests
	go test ./...; \
	find . -type d -name "logs" -exec rm -r {} \;

doc: ## Generate API documentation using swag
	swag init -g cmd/api/main.go

<<<<<<< HEAD
prod: ##push prod
=======
prod: lint doc
>>>>>>> 493d329e5b644d4e30dec179c1f48f05106223bb
	git checkout deploy ; \
	git pull origin develop ; \
	git add . ; \
	git commit -m "deploy" ; \
<<<<<<< HEAD
	git push
=======
	git push ; \
	git checkout develop
>>>>>>> 493d329e5b644d4e30dec179c1f48f05106223bb
