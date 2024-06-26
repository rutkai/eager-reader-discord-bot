FROM golang:1.22-bullseye as builder

WORKDIR /app

COPY src/go.* ./
RUN go mod download

COPY src/ ./
RUN go build -v -o eager-reader-discord-bot

FROM debian:bullseye-slim

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/eager-reader-discord-bot /app/eager-reader-discord-bot

ENTRYPOINT ["/app/eager-reader-discord-bot"]
