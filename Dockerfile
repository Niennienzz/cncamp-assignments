# syntax=docker/dockerfile:1

### Build Stage ###
# Choose a good enough base image for building: alpine.
FROM golang:1.16-alpine AS build-env

# SQLite is an embedded C library.
# The HTTP server uses SQLite, and thus requires the gcc compiler.
RUN set -ex && apk add --no-cache gcc musl-dev

# Copy source files.
WORKDIR /app
COPY . .

# Resolve dependencies.
RUN go mod vendor

# Compile to binary.
RUN go build -o /cncamp_http_server

### Final Stage ###
# Choose one of the smallest base image: distroless.
FROM gcr.io/distroless/base-debian10

COPY --from=build-env /cncamp_http_server /

EXPOSE 8080
CMD [ "/cncamp_http_server" ]
