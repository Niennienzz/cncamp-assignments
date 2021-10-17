# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

# SQLite is embedded and thus github.com/mattn/go-sqlite3 requires the gcc compiler.
RUN set -ex && \
    apk add --no-cache gcc musl-dev

WORKDIR /app
COPY . .

RUN go mod vendor
RUN go build -o /cncamp_http_server

EXPOSE 8080
CMD [ "/cncamp_http_server" ]
