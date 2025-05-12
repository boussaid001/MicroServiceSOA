FROM golang:1.20-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git protobuf-dev

# Install protoc-gen-go and protoc-gen-go-grpc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Copy proto files first
COPY proto/ /app/proto/

# Regenerate proto files
RUN mkdir -p /app/proto/generated
RUN protoc --go_out=/app/proto/generated --go-grpc_out=/app/proto/generated -I=/app/proto /app/proto/product.proto

# Copy go.mod and go.sum first to leverage Docker cache
COPY services/grpc-service/go.mod services/grpc-service/go.sum* services/grpc-service/

# Install Go modules
RUN cd services/grpc-service && go mod download && go mod tidy

# Copy the rest of the service code
COPY services/grpc-service/ services/grpc-service/

# Build the gRPC service
RUN cd services/grpc-service && go build -o /app/main cmd/main.go

# Expose the gRPC port
EXPOSE 8082

# Run the service
CMD ["/app/main"] 