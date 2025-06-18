# This docker file is the base for each GO serviuce
FROM golang:1.24

# Set working directory
WORKDIR /app

# Copy your project files into the container
COPY . .

# Install dependencies
RUN go mod download
