# Base image for Go
FROM golang:1.24-alpine AS builder

# Working directory inside the container
WORKDIR /app

RUN apk add --no-cache musl-dev gcc

# Copy the Go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application (main.go is in the project root); CGO has to be enabled because of valkey-glide
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64 GOFLAGS=-tags=musl
RUN go build -o hidden-live-location-server .

# Runtime image
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache musl-dev gcc

# Copy the built Go binary
COPY --from=builder /app/hidden-live-location-server .

EXPOSE 8080

CMD ["./hidden-live-location-server"]
