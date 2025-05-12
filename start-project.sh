#!/bin/bash

echo "Starting Go Microservices Project..."

# Define the ports used by the project
PORTS=(8080 8081 8082 8083 29092 32181)

# Stop any existing docker-compose services first
echo "Stopping any existing containers..."
docker-compose down 2>/dev/null

# Specifically check for bettercap on port 8081 (which keeps reappearing)
echo "Checking specifically for bettercap on port 8081..."
BETTERCAP_PID=$(sudo lsof -i :8081 | grep bettercap | awk '{print $2}')
if [ ! -z "$BETTERCAP_PID" ]; then
  echo "Found bettercap (PID: $BETTERCAP_PID) on port 8081. Killing process..."
  sudo kill -9 $BETTERCAP_PID
  echo "Process killed."
fi

# Check if any ports are in use and kill those processes (with sudo)
for PORT in "${PORTS[@]}"; do
  # First try normal lsof
  PID=$(lsof -i :$PORT -t 2>/dev/null)
  if [ ! -z "$PID" ]; then
    echo "Port $PORT is in use by process $PID. Killing process..."
    kill -9 $PID 2>/dev/null || sudo kill -9 $PID
    echo "Process killed."
  else
    # If normal lsof didn't find anything, try with sudo
    echo "Checking for processes using port $PORT that might require sudo..."
    PID=$(sudo lsof -i :$PORT -t 2>/dev/null)
    if [ ! -z "$PID" ]; then
      echo "Port $PORT is in use by process $PID (requires sudo). Killing process..."
      sudo kill -9 $PID
      echo "Process killed."
    fi
  fi
done

# Additional check for Docker port bindings
echo "Checking for Docker containers using our ports..."
for PORT in "${PORTS[@]}"; do
  CONTAINER=$(docker ps | grep -E ":[0-9]+->$PORT" | awk '{print $1}')
  if [ ! -z "$CONTAINER" ]; then
    echo "Docker container $CONTAINER is using port $PORT. Stopping container..."
    docker stop $CONTAINER
    echo "Container stopped."
  fi
done

# Make sure Docker is running
if ! docker info > /dev/null 2>&1; then
  echo "Docker is not running. Starting docker service..."
  sudo systemctl start docker
  if [ $? -ne 0 ]; then
    echo "Failed to start Docker. Please start it manually."
    exit 1
  fi
fi

# Start the project with Docker Compose
echo "Starting services with Docker Compose..."
docker-compose up -d

# Check if all services started successfully
if [ $? -eq 0 ]; then
  echo "Project started! Services should be available at the following endpoints:"
  echo "API Gateway:        http://localhost:8080"
  echo "REST Service:       http://localhost:8081"
  echo "gRPC Service:       http://localhost:8082 (for client connections)"
  echo "GraphQL Service:    http://localhost:8083"
  echo "Kafka:              localhost:29092"
  echo "Zookeeper:          localhost:32181"
  echo
  echo "To view logs, run: docker-compose logs -f"
  echo "To stop services, run: docker-compose down"
else
  echo "Some services failed to start. Check logs for details: docker-compose logs"
fi 