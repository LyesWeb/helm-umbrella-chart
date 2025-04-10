FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod init countries-service && \
    go build -o countries-service

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/countries-service /app/
COPY countries.csv /app/

EXPOSE 8080
HEALTHCHECK --interval=5s --timeout=3s --retries=3 CMD wget -qO- http://localhost:8080/health || exit 1

CMD ["./countries-service"]