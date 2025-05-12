#!/bin/bash

echo "Starting Go Microservices Project..."

# Define color codes for better readability
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Define the ports used by the project
PORTS=(8080 8081 8082 8083 8090 29092 32181 5432)

# Stop any existing docker-compose services first
echo "Stopping any existing containers..."
docker-compose down 2>/dev/null

# Function to check if a port is in use
check_and_free_port() {
  local PORT=$1
  # Check for processes using the port
  PID=$(lsof -i :$PORT -t 2>/dev/null)
  if [ ! -z "$PID" ]; then
    echo -e "${YELLOW}Port $PORT is in use by process $PID. Killing process...${NC}"
    kill -9 $PID 2>/dev/null || sudo kill -9 $PID
    echo "Process killed."
  else
    # If normal lsof didn't find anything, try with sudo
    PID=$(sudo lsof -i :$PORT -t 2>/dev/null 2>&1)
    if [ ! -z "$PID" ]; then
      echo -e "${YELLOW}Port $PORT is in use by process $PID (requires sudo). Killing process...${NC}"
      sudo kill -9 $PID
      echo "Process killed."
    fi
  fi
}

# Free all required ports
echo "Checking for processes using required ports..."
for PORT in "${PORTS[@]}"; do
  check_and_free_port $PORT
done

# Make sure Docker is running
if ! docker info > /dev/null 2>&1; then
  echo -e "${YELLOW}Docker is not running. Starting docker service...${NC}"
  sudo systemctl start docker
  if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to start Docker. Please start it manually.${NC}"
    exit 1
  fi
fi

# Ask the user if they want to rebuild containers
echo -n -e "${YELLOW}Do you want to rebuild the containers? (y/n): ${NC}"
read -r REBUILD
if [[ "$REBUILD" =~ ^[Yy]$ ]]; then
  echo "Rebuilding containers..."
  docker-compose build --no-cache api-gateway grpc-service
  if [ $? -ne 0 ]; then
    echo -e "${RED}Failed to rebuild containers. Proceeding with existing containers.${NC}"
  else
    echo -e "${GREEN}Containers rebuilt successfully.${NC}"
  fi
fi

# Start the project with Docker Compose
echo -e "${GREEN}Starting services with Docker Compose...${NC}"
docker-compose up -d

# Check if all services started successfully
if [ $? -eq 0 ]; then
  echo -e "${GREEN}All services started!${NC}"
  
  # Wait for services to be fully initialized
  echo "Waiting for services to initialize..."
  sleep 10
  
  # Check health of API Gateway
  echo "Checking API Gateway health..."
  if curl -s http://localhost:8080/health | grep -q "ok"; then
    echo -e "${GREEN}✅ API Gateway is healthy${NC}"
  else
    echo -e "${RED}❌ API Gateway might not be fully initialized${NC}"
  fi
  
  # Check gRPC Service by listing products
  echo "Checking gRPC Service..."
  if curl -s http://localhost:8080/api/products/ | grep -q "products"; then
    echo -e "${GREEN}✅ gRPC Service is responsive${NC}"
  else
    echo -e "${RED}❌ gRPC Service might not be fully initialized${NC}"
  fi
  
  echo -e "\n${GREEN}Services should be available at the following endpoints:${NC}"
  echo "API Gateway:        http://localhost:8080"
  echo "Frontend:           http://localhost:8080/products.html"
  echo "REST Service:       http://localhost:8081"
  echo "gRPC Service:       http://localhost:8082 (for client connections)"
  echo "GraphQL Service:    http://localhost:8083"
  echo "Hasura Console:     http://localhost:8090"
  echo "Kafka:              localhost:29092"
  echo "Zookeeper:          localhost:32181"
  echo
  echo -e "${YELLOW}To view logs, run: docker-compose logs -f${NC}"
  echo -e "${YELLOW}To stop services, run: docker-compose down${NC}"
else
  echo -e "${RED}Some services failed to start. Check logs for details: docker-compose logs${NC}"
fi 