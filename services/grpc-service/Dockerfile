FROM golang:1.21-alpine as builder
WORKDIR /app

# Install git and protobuf compiler
RUN apk add --no-cache git protobuf protobuf-dev

# Copy proto files
COPY ./proto /app/proto

# Copy service code
COPY ./services/grpc-service /app/services/grpc-service
WORKDIR /app/services/grpc-service

# Initialize and set up Go module with dependencies
RUN go mod init github.com/yourusername/go-microservices-project/services/grpc-service
RUN go mod edit -require=google.golang.org/grpc@v1.58.3
RUN go mod edit -replace=github.com/yourusername/go-microservices-project/proto/product=../../proto/product

# Build the application
RUN go mod tidy
RUN go build -o main cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/services/grpc-service/main .
EXPOSE 8082
CMD ["./main"] 