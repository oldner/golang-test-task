FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o reporting-api ./cmd/reporting-api

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/reporting-api .
COPY .env .

EXPOSE 8081

CMD ["./reporting-api"]