# Basis-Image für Go
FROM golang:1.24 AS builder

# Arbeitsverzeichnis im Container
WORKDIR /app

# Kopiere die Go-Moduldateien und source code
COPY ../go.mod ../go.sum ./
RUN go mod download

# Kopiere alle Go-Dateien in das Arbeitsverzeichnis
COPY ../*.go ./

# Baue die Go-Anwendung
RUN go build -o pps

# Laufzeit-Image
FROM alpine:latest

# Arbeitsverzeichnis im Container
WORKDIR /app

# Kopiere das gebaute Go-Programm
COPY --from=builder /app/pps .

# Exponiere den Port, auf dem die App läuft
EXPOSE 8080

# Starte die Go-Anwendung
CMD ["./pps"]
