FROM golang:1.24

WORKDIR /app
COPY . .

# Install dependencies
RUN go mod download

# Install specific versions of dependencies (optional)
# RUN go install github.com/lib/pq@v1.10.9
# RUN go install github.com/pressly/goose/v3/cmd/goose@v3.24.3

# Default beheviour
# ENTRYPOINT ["go", "run", "main.go", "up", "-migrations",  "./migrations"]
