# Go Microservices Project

This project demonstrates a microservices architecture using Go, featuring an API Gateway, gRPC, REST, GraphQL, and Kafka integration.

## Project Structure

- **api-gateway/**  
  Contains the API Gateway service that routes requests to the appropriate microservice.  
  - `cmd/main.go`: Entry point for the API Gateway service.  
  - `internal/`: Internal packages for the API Gateway.  
  - `Dockerfile`: Dockerfile for building the API Gateway container.  
  - `go.mod` & `go.sum`: Go module files for dependency management.

- **grpc-service/**  
  Contains the gRPC service for handling product operations.  
  - `cmd/main.go`: Entry point for the gRPC service.  
  - `internal/`: Internal packages for the gRPC service.  
  - `Dockerfile`: Dockerfile for building the gRPC service container.  
  - `go.mod` & `go.sum`: Go module files for dependency management.

- **rest-service/**  
  Contains the REST service for handling product operations.  
  - `cmd/main.go`: Entry point for the REST service.  
  - `internal/`: Internal packages for the REST service.  
  - `Dockerfile`: Dockerfile for building the REST service container.  
  - `go.mod` & `go.sum`: Go module files for dependency management.

- **graphql-service/**  
  Contains the GraphQL service for handling product operations.  
  - `cmd/main.go`: Entry point for the GraphQL service.  
  - `internal/`: Internal packages for the GraphQL service.  
  - `Dockerfile`: Dockerfile for building the GraphQL service container.  
  - `go.mod` & `go.sum`: Go module files for dependency management.

- **kafka-service/**  
  Contains the Kafka service for event processing.  
  - `cmd/main.go`: Entry point for the Kafka service.  
  - `internal/`: Internal packages for the Kafka service.  
  - `Dockerfile`: Dockerfile for building the Kafka service container.  
  - `go.mod` & `go.sum`: Go module files for dependency management.

- **frontend/**  
  Contains the frontend application for interacting with the microservices.  
  - `index.html`: Main HTML file for the frontend.  
  - `products.html`: HTML file for the products page.  
  - `styles.css`: CSS file for styling the frontend.  
  - `script.js`: JavaScript file for frontend logic.

- **docker-compose.yml**  
  Defines the services, networks, and volumes for the project. It orchestrates the deployment of all microservices, including the API Gateway, gRPC, REST, GraphQL, Kafka, and Hasura.

- **start-project.sh**  
  A shell script to start the project. It stops any existing containers, checks for processes using required ports, and starts the services using Docker Compose.

- **postman-collection.json**  
  A Postman collection for testing all API endpoints of the project.

- **postman-environment.json**  
  A Postman environment file for setting up environment variables for local development.

- **README-postman.md**  
  Documentation on how to use the Postman collection, including import instructions, environment setup, and API references.

- **README.md**  
  This file, providing an overview of the project, its structure, and instructions for running the project.

## Running the Project

1. **Start the Project:**  
   Run the following command to start all services:  
   ```bash
   ./start-project.sh
   ```

2. **Access the Services:**  
   - API Gateway: [http://localhost:8080](http://localhost:8080)
   - Frontend: [http://localhost:8080/products.html](http://localhost:8080/products.html)
   - REST Service: [http://localhost:8081](http://localhost:8081)
   - gRPC Service: [http://localhost:8082](http://localhost:8082) (for client connections)
   - GraphQL Service: [http://localhost:8083](http://localhost:8083)
   - Hasura Console: [http://localhost:8090](http://localhost:8090)
   - Kafka: localhost:29092
   - Zookeeper: localhost:32181

3. **Testing the APIs:**  
   Import the Postman collection (`postman-collection.json`) and environment (`postman-environment.json`) into Postman to test all API endpoints.

4. **Viewing Logs:**  
   To view logs, run:  
   ```bash
   docker-compose logs -f
   ```

5. **Stopping the Services:**  
   To stop all services, run:  
   ```bash
   docker-compose down
   ```

## Conclusion

This project demonstrates a microservices architecture using Go, featuring an API Gateway, gRPC, REST, GraphQL, and Kafka integration. Each microservice is containerized using Docker, and the project is orchestrated using Docker Compose. The frontend provides a user interface for interacting with the microservices, and the Postman collection provides a comprehensive set of tests for all API endpoints.
