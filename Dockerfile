# Specify the Go version
ARG GO_VERSION=1.24.3

# Use a Go image with the specified version for the build stage
FROM golang:${GO_VERSION}-alpine AS build

ENV CGO_ENABLED=0

WORKDIR /app

COPY . .

ENV GOFLAGS="-buildmode=pie"

RUN apk update --no-cache ca-certificates
RUN go mod download
RUN go build -o /app/main ./cmd/watcher

# Start from scratch for the final image
FROM scratch AS app

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/main main

# Define the command to run the application
CMD ["/main"]