FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o message-processor ./cmd/message-processor

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/message-processor .
COPY .env .

CMD ["./message-processor"]