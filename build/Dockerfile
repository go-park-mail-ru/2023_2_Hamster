#Builder
FROM golang:1.21.0-alpine AS builder

COPY . /github.com/go-park-mail-ru/2023_2_Hamster/
WORKDIR /github.com/go-park-mail-ru/2023_2_Hamster/

RUN go mod download
RUN go clean --modcache
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o server ./cmd/api/main.go

FROM golang:latest AS run

WORKDIR /docker-hammywallet/

COPY --from=builder /github.com/go-park-mail-ru/2023_2_Hamster/server .

EXPOSE 80 433 8080

ENTRYPOINT ["./server"]
