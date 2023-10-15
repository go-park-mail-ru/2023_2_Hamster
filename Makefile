include .env

.PHONY: dc run test lint down

run:
	docker-compose up -d
build:
	docker-compose build

сondb:
	docker exec -it 2023_2_hamster-db-1 psql -U kosmatoff -d Hamster
down:
	docker-compose down

doc:
	swag init -g cmd/api/main.go