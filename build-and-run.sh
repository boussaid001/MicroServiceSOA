#!/bin/bash

# Exit on error
set -e

echo "Building and running Go Microservices with dual GraphQL backends..."

# Clean up any existing containers to avoid conflicts
echo "Cleaning up old containers..."
docker-compose down

# Build fresh containers
echo "Building containers..."
docker-compose build

# Start the system
echo "Starting the system..."
docker-compose up -d

# Wait for services to be ready
echo "Waiting for services to initialize..."
sleep 10

# Show running containers
echo "Running containers:"
docker-compose ps

echo "Setting up Hasura..."
# Create tables for Hasura if needed
# Initialize permissions

echo "System is ready! Access it at:"
echo "- API Gateway: http://localhost:8080"
echo "- Custom GraphQL: http://localhost:8083/graphql"
echo "- Hasura Console: http://localhost:8090/console"
echo "- GraphQL Comparison UI: http://localhost:8080/compare-graphql.html"
echo ""
echo "View logs with: docker-compose logs -f" 