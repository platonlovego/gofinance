FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o gofinance ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/gofinance .

EXPOSE 8080
CMD ["./gofinance"]
