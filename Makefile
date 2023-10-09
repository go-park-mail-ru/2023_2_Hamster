include .env

.PHONY: dc run test lint down

dc:
	sudo docker-compose up -d

run:
	go build -o app ./cmd/api/main.go && ./app

down:
	sudo docker-compose down