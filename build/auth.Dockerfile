#Builder
FROM golang:1.21.0-alpine AS builder

COPY . /github.com/go-park-mail-ru/2023_2_Hamster/
WORKDIR /github.com/go-park-mail-ru/2023_2_Hamster/

RUN go mod download && go mod tidy
RUN go clean --modcache
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o auth ./cmd/auth/main.go

FROM golang:1.21.0-alpine AS run

WORKDIR /docker-hammywallet/

COPY --from=builder /github.com/go-park-mail-ru/2023_2_Hamster/auth .

EXPOSE 8010

ENTRYPOINT ["./auth"]
