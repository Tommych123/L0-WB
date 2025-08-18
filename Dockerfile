FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/l0-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/l0-service /app/l0-service

COPY ./web /app/web

COPY .env /app/.env

EXPOSE 8080

CMD ["/app/l0-service"]