FROM golang:1.21-alpine as builder
WORKDIR /app

# Install git
RUN apk add --no-cache git

# Copy proto files 
COPY ./proto /app/proto

# Copy API Gateway code
COPY ./api-gateway /app/api-gateway
WORKDIR /app/api-gateway

# Initialize and set up Go module with dependencies
RUN go mod init github.com/boussaid001/go-microservices-project/api-gateway
RUN go mod edit -require=golang.org/x/crypto@v0.16.0
RUN go mod edit -require=github.com/IBM/sarama@v1.42.1
RUN go mod edit -require=google.golang.org/grpc@v1.58.3
RUN go mod edit -require=github.com/gin-gonic/gin@v1.9.1
RUN go mod edit -require=github.com/rogpeppe/go-internal@v1.11.0
RUN go mod edit -require=github.com/gin-contrib/cors@v1.4.0
RUN go mod edit -replace=github.com/boussaid001/go-microservices-project/proto=../proto

# Build the application
RUN go mod tidy
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api-gateway/main .

# Copy frontend files
COPY ./frontend /app/frontend

EXPOSE 8080
CMD ["./main"] 