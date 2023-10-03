# Используйте официальный образ Golang как базовый образ
FROM golang:latest

# Установите рабочую директорию внутри контейнера
WORKDIR /2023_2_Hamster

# Скопируйте файлы вашего Golang приложения в контейнер
COPY . .


RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x wait-for-postgres.sh

# Соберите ваше приложение (замените на имя вашего исполняемого файла, если не "main")
RUN go mod download
RUN go build -o app cmd/api/main.go


# Укажите команду для запуска вашего приложения при старте контейнера
CMD ["./app"]