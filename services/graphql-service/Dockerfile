# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy the source code
COPY ./services/graphql-service /app

# Initialize go modules
RUN go mod init github.com/yourusername/go-microservices-project/services/graphql-service
RUN go mod edit -require=github.com/graph-gophers/graphql-go@v1.5.0
RUN go mod edit -require=github.com/lib/pq@v1.10.9
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o graphql-service ./cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary and schema files
COPY --from=builder /app/graphql-service .
COPY --from=builder /app/schema /app/schema

# Expose port
EXPOSE 8083

# Run the application
CMD ["./graphql-service"] 