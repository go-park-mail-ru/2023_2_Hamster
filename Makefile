include .env

.PHONY: dc run test lint down

run:
	docker-compose up -d

down:
	docker-compose down

doc:
	swag init -g cmd/api/main.go