include .env

.PHONY: dc run test lint down

run:
	docker-compose up -d

run-in:
	docker-compose up

dock-build:
	docker-compose build

dock-clean:
	docker rmi -f $(docker images -q)

con-db:
	docker exec -it 2023_2_hamster-db-1 psql -U kosmatoff -d Hamster

con-app:
	docker exec -it 2023_2_hamster-server-1 ./app

down:
	docker-compose down

doc:
	swag init -g cmd/api/main.go