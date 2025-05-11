FROM golang:1.21-alpine as builder
WORKDIR /app

# Install git
RUN apk add --no-cache git

# Copy service code
COPY ./services/kafka-service /app

# Initialize and set up Go module with dependencies
RUN go mod init github.com/yourusername/go-microservices-project/services/kafka-service
RUN go mod edit -require=github.com/IBM/sarama@v1.42.1

# Build the application
RUN go mod tidy
RUN go build -o main cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8084
CMD ["./main"] 