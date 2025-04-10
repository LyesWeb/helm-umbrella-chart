FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod init client-service && \
    go build -o client-service

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/client-service /app/

EXPOSE 8081
HEALTHCHECK --interval=5s --timeout=3s --retries=3 CMD wget -qO- http://localhost:8081/health || exit 1

CMD ["./client-service"]