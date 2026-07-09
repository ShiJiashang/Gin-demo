FROM golang:1.26.3-bookworm AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/gin-demo \
    .

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /out/gin-demo /app/gin-demo

RUN mkdir -p /data

ENV PORT=8080
ENV DB_PATH=/data/gin-demo.db

EXPOSE 8080

CMD ["/app/gin-demo"]
