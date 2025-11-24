# Build stage
FROM golang:1.23-alpine AS builder

# Install git (used by some Go modules)
RUN apk add --no-cache git

# Application workspace
WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the binary in the same way as the working deployment
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./main.go

# Final stage
FROM alpine:latest

# Only ca-certificates are required at runtime
RUN apk --no-cache add ca-certificates

# Match the working deployment layout
WORKDIR /root/

# Copy the compiled binary
COPY --from=builder /app/main .

# Expose service port
EXPOSE 8080

# Run the service
CMD ["./main"]

