FROM golang:1.23.4

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

COPY config/config.yaml ./config.yaml
COPY migrations ./migrations

RUN go build -o app ./cmd/app

CMD ["./app"]
