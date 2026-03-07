# Base image for Go
FROM golang:1.24 AS builder

# Working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application (main.go is in the project root)
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o pps .

# Runtime image
FROM alpine:latest

WORKDIR /app

# Copy the built Go binary
COPY --from=builder /app/pps .

EXPOSE 8080

CMD ["./pps"]
