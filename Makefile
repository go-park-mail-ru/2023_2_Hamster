include .env

.PHONY: dc run test lint down

dc:
	docker-compose up -d

run:
	go build -o app ./cmd/api/main.go && ./app

down:
	docker-compose down