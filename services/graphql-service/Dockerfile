FROM golang:1.21-alpine as builder
WORKDIR /app

# Install git
RUN apk add --no-cache git

# Copy service code
COPY ./services/graphql-service /app

# Initialize and set up Go module with dependencies
RUN go mod init github.com/yourusername/go-microservices-project/services/graphql-service
RUN go mod edit -require=github.com/graphql-go/graphql@v0.8.1
RUN go mod edit -require=github.com/graphql-go/handler@v0.2.3
RUN go mod edit -require=github.com/gin-gonic/gin@v1.9.1
RUN go mod edit -require=github.com/rogpeppe/go-internal@v1.11.0

# Build the application
RUN go mod tidy
RUN go build -o main cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8083
CMD ["./main"] 