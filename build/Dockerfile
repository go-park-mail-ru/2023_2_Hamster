FROM golang:1.21.0-alpine AS build

COPY . /github.com/go-park-mail-ru/2023_2_Hamster/
WORKDIR /github.com/go-park-mail-ru/2023_2_Hamster/

RUN go mod download
RUN go clean --modcache
RUN go build -mod=readonly -o ./.server ./cmd/api/main.go

FROM scratch AS run

WORKDIR /docker-HammyWallet/

COPY --from=builder /github.com/go-park-mail-ru/2023_2_Hamster/.server .

EXPOSE 80 433

ENTRYPOINT ["./server"]