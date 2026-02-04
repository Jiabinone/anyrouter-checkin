FROM node:20-bookworm-slim AS frontend-build

WORKDIR /src/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

FROM golang:1.23-bookworm AS backend-build

WORKDIR /src/backend
RUN apt-get update \
  && apt-get install -y --no-install-recommends build-essential pkg-config libsqlite3-dev \
  && rm -rf /var/lib/apt/lists/*
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=1 go build -o /bin/server cmd/server/main.go

FROM debian:bookworm-slim

ENV TZ=Asia/Shanghai
ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates sqlite3 libsqlite3-0 nginx supervisor tzdata \
  && ln -snf /usr/share/zoneinfo/${TZ} /etc/localtime \
  && echo "${TZ}" > /etc/timezone \
  && rm -f /etc/nginx/sites-enabled/default /etc/nginx/sites-available/default \
  && rm -rf /var/lib/apt/lists/* \
  && mkdir -p /run/nginx /var/log/supervisor /app/data

WORKDIR /app

COPY --from=backend-build /bin/server /app/server
COPY backend/config.example.yaml /app/config.yaml
COPY --from=frontend-build /src/frontend/dist /usr/share/nginx/html
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf
COPY docker/supervisord.conf /etc/supervisord.conf

EXPOSE 80 8080

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
