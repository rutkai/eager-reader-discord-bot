FROM golang:1.22-buster as builder

WORKDIR /app

COPY src/go.* ./
RUN go mod download

COPY src/ ./
RUN go build -v -o eager-reader-discord-bot

FROM debian:buster-slim

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/eager-reader-discord-bot /app/eager-reader-discord-bot

CMD ["/app/eager-reader-discord-bot"]
