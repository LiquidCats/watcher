# Specify the Go version
ARG GO_VERSION=1.24

# Use a Go image with the specified version for the build stage
FROM golang:${GO_VERSION}-alpine AS build

ENV CGO_ENABLED=0

WORKDIR /app

COPY . .

ENV GOFLAGS="-buildmode=pie"

RUN go mod download
RUN go build -o /app/main ./cmd/watcher

# Start from scratch for the final image
FROM scratch AS app

WORKDIR /

USER 65534

# Copy the built binary from the build stage
COPY --from=build /app/main main

# Define the command to run the application
CMD ["/main"]