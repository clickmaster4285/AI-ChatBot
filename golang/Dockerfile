FROM golang:1.22-alpine

WORKDIR /app

# Copy everything first (safe, flexible)
COPY . .

# Download deps only if they exist
RUN go mod tidy || true

# Build binary
RUN go build -o main cmd/server/main.go

EXPOSE 4000

CMD ["./main"]