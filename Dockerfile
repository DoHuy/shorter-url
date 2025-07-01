# Build stage
FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o shorter-rest-api .

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/shorter-rest-api .

EXPOSE 8080

CMD ["./shorter-rest-api"]