FROM golang:latest
COPY . /app
WORKDIR /app
RUN go mod download
RUN go build -o app ./cmd/api/main.go
ENTRYPOINT ["./app"]
