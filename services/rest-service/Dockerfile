FROM golang:1.21-alpine as builder
WORKDIR /app

# Install git
RUN apk add --no-cache git

# Copy service code
COPY ./services/rest-service /app

# Initialize and set up Go module with dependencies
RUN go mod init github.com/yourusername/go-microservices-project/services/rest-service
RUN go mod edit -require=golang.org/x/crypto@v0.16.0
RUN go mod edit -require=github.com/gin-gonic/gin@v1.9.1
RUN go mod edit -require=github.com/rogpeppe/go-internal@v1.11.0

# Build the application
RUN go mod tidy
RUN go build -o main cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8081
CMD ["./main"] 