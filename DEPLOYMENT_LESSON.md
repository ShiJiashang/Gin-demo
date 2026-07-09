# Gin + Docker + Nginx deployment lesson

## Architecture

```text
                         +-> gin-demo-1 (host port 8081)
client -> Nginx :9090 ---|
                         +-> gin-demo-2 (host port 8082)
```

Nginx listens on port `9090`. Its `upstream gin_backend` group contains two
instances of the same Docker image. With the default round-robin strategy,
Nginx selects the next instance for each request.

## Build the image

```bash
docker build -t gin-demo:lesson .
```

The `Dockerfile` has two stages:

1. `golang:1.26.3-bookworm` downloads dependencies and compiles the program.
2. `debian:bookworm-slim` contains only the binary and runtime libraries.

## Start two Gin containers

```bash
docker run -d \
  --name gin-demo-1 \
  -p 127.0.0.1:8081:8080 \
  -e JWT_SECRET=lesson-secret \
  -v gin-demo-1-data:/data \
  gin-demo:lesson

docker run -d \
  --name gin-demo-2 \
  -p 127.0.0.1:8082:8080 \
  -e JWT_SECRET=lesson-secret \
  -v gin-demo-2-data:/data \
  gin-demo:lesson
```

`8081:8080` means host port `8081` is forwarded to container port `8080`.
Each container uses its own named volume so deleting the container does not
delete its SQLite database.

These two SQLite databases are intentionally independent for this deployment
exercise. Real replicas must share an external database such as MySQL or
PostgreSQL; otherwise writes sent to different instances produce inconsistent
data.

## Check and reload Nginx

Run these commands from the project root:

```bash
nginx -p "$PWD/nginx-demo" -c conf/nginx.conf -t
nginx -p "$PWD/nginx-demo" -c conf/nginx.conf -s reload
```

If Nginx is not running yet:

```bash
nginx -p "$PWD/nginx-demo" -c conf/nginx.conf
```

## Verify round-robin load balancing

```bash
for i in 1 2 3 4; do
  curl -s -o /dev/null -D - http://localhost:9090/swagger/doc.json \
    | grep -i x-upstream-addr
done
```

Expected response headers alternate:

```text
X-Upstream-Addr: 127.0.0.1:8081
X-Upstream-Addr: 127.0.0.1:8082
```

## Stop and remove the lesson containers

```bash
docker rm -f gin-demo-1 gin-demo-2
```

The named volumes remain. Delete them only when their SQLite data is no longer
needed:

```bash
docker volume rm gin-demo-1-data gin-demo-2-data
```
