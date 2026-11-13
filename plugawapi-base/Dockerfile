FROM golang:1.26 AS builder
WORKDIR /app
COPY . .
RUN go build -o plugawapi ./cmd/server

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/plugawapi .
COPY --from=builder /app/client/dist ./client/dist
EXPOSE 8080
CMD ["./plugawapi", "-addr", ":8080", "-static", "client/dist"]
