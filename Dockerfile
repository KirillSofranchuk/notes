FROM golang:1.23.4

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Копируем исходный код проекта
COPY . ./

COPY config/config.yaml ./config.yaml
COPY migrations ./migrations

# Сборка бинарника (в указанной директории)
RUN go build -o app ./cmd/app

CMD ["./app"]
