# Go Microservices Project

A comprehensive microservices-based application showcasing different communication patterns and technologies, built using Go and Docker.

## Author
**Mohamed Amine Boussaid**

## Table of Contents
- [Architecture Overview](#architecture-overview)
- [Services](#services)
  - [API Gateway](#api-gateway)
  - [REST User Service](#rest-user-service)
  - [gRPC Product Service](#grpc-product-service)
  - [Hasura GraphQL (Reviews)](#hasura-graphql-reviews)
  - [Kafka Order Service](#kafka-order-service)
- [Database Structure](#database-structure)
- [Communication Patterns](#communication-patterns)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Technology Stack](#technology-stack)

## Architecture Overview

```
┌─────────────┐     ┌─────────────────────────────────────────────────┐
│             │     │                                                 │
│  Frontend   │◄────┤                    API Gateway                  │
│  (Browser)  │     │                  (Port: 8080)                   │
│             │     │                                                 │
└──────┬──────┘     └─┬───────────────────┬───────────────┬──────────┘
       │              │                   │               │
       │              │                   │               │
       ▼              ▼                   ▼               ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐ ┌─────────────┐     ┌─────────────┐
│             │     │             │     │             │ │             │     │             │
│   Browser   │     │  REST API   │     │ gRPC Service│ │   Hasura    │◄───┤Kafka Service │
│  Direct to  │     │ User Service│     │   Product   │ │   GraphQL   │     │   Orders    │
│   Hasura    │     │ (Port:8081) │     │ (Port:8082) │ │ (Port:8090) │     │             │
│             │     │             │     │             │ │             │     │             │
└──────┬──────┘     └──────┬──────┘     └──────┬──────┘ └──────┬──────┘     └──────┬──────┘
       │                   │                   │               │                   │
       │                   │                   │               │                   │
       ▼                   ▼                   ▼               ▼                   ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│             │     │             │     │             │     │             │     │             │
│  PostgreSQL │     │  PostgreSQL │     │  PostgreSQL │     │  PostgreSQL │     │    Kafka    │
│  reviewdb   │     │   userdb    │     │  productdb  │     │   orderdb   │     │   Broker    │
│             │     │             │     │             │     │             │     │             │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
```

This project demonstrates a microservices architecture with services communicating via different protocols:
- **REST API**: Used by the User Service
- **gRPC**: Used by the Product Service
- **GraphQL**: Used by Hasura for Reviews 
- **Kafka**: Used for event-driven communication in the Order Service

## Services

### API Gateway
- **Port**: 8080
- **Role**: Entry point for all client requests
- **Responsibilities**:
  - Route requests to appropriate backend services
  - Serve static frontend files
  - Proxy GraphQL requests to Hasura
  - Handle CORS and authentication
- **Technologies**: Go, Gin Framework
- **Communication**: HTTP/REST with clients and backend services

### REST User Service
- **Port**: 8081
- **Role**: Manage user data and authentication
- **Endpoints**:
  - `GET /api/users`: List all users
  - `GET /api/users/:id`: Get user by ID
  - `POST /api/users`: Create a new user
  - `PUT /api/users/:id`: Update a user
  - `DELETE /api/users/:id`: Delete a user
- **Technologies**: Go, Gin Framework
- **Storage**: PostgreSQL (userdb)
- **Communication**: REST/HTTP

### gRPC Product Service
- **Port**: 8082
- **Role**: Manage product catalog
- **Operations**:
  - Get all products
  - Get product by ID
  - Create product
  - Update product
  - Delete product
- **Technologies**: Go, gRPC, Protocol Buffers
- **Storage**: PostgreSQL (productdb)
- **Communication**: gRPC (binary protocol over HTTP/2)

### Hasura GraphQL (Reviews)
- **Port**: 8090
- **Role**: Manage product reviews
- **Operations**:
  - Query reviews (by ID, by product)
  - Create reviews
- **Technologies**: Hasura GraphQL Engine
- **Storage**: PostgreSQL (reviewdb)
- **Communication**: GraphQL over HTTP
- **Schema**:
  ```graphql
  type Review {
    id: String!
    product_id: String!
    user_id: String!
    username: String!
    rating: Float!
    comment: String
    created_at: timestamptz!
  }
  ```

### Kafka Order Service
- **Role**: Process and track orders
- **Events**:
  - Order Created
  - Order Updated
  - Order Status Updated
- **Technologies**: Go, Kafka
- **Storage**: PostgreSQL (orderdb)
- **Communication**: Event-driven messaging via Kafka

## Database Structure

Each service has its own dedicated PostgreSQL database:

### UserDB
- **Tables**:
  - `users`: Stores user information (id, name, email, password, created_at)

### ProductDB
- **Tables**:
  - `products`: Stores product information (id, name, description, price, stock, created_at)

### ReviewDB
- **Tables**:
  - `reviews`: Stores review information (id, product_id, user_id, username, rating, comment, created_at)

### OrderDB
- **Tables**:
  - `orders`: Stores order information (id, user_id, status, total, created_at)
  - `order_items`: Stores items in each order (order_id, product_id, quantity, price)

## Communication Patterns

This project demonstrates multiple communication patterns:

1. **Synchronous Communication**:
   - REST API (User Service)
   - gRPC (Product Service)
   - GraphQL (Reviews via Hasura)

2. **Asynchronous Communication**:
   - Kafka messaging (Order Service)

3. **Client-Service Communication**:
   - Browser to API Gateway (HTTP/REST)
   - Browser to Hasura (Direct GraphQL)

4. **Service-to-Service Communication**:
   - API Gateway to backend services

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Go 1.21 or later (for development)

### Running the Application
1. Clone the repository:
   ```
   git clone https://github.com/yourusername/go-microservices-project.git
   cd go-microservices-project
   ```

2. Start the services:
   ```
   ./start-project.sh
   ```

3. Access the application:
   - Dashboard: http://localhost:8080
   - Users page: http://localhost:8080/users.html
   - Products page: http://localhost:8080/products.html
   - Reviews page: http://localhost:8080/reviews.html
   - Orders page: http://localhost:8080/orders.html
   - Hasura Console: http://localhost:8090

### Stopping the Application
```
docker-compose down
```

## API Documentation

API documentation is available in the `Docs.json` file, which is a Postman collection that can be imported to test all APIs.

## Technology Stack

- **Languages**: Go, JavaScript (frontend)
- **API Gateway**: Go with Gin framework
- **Services**:
  - REST: Go with Gin framework
  - gRPC: Go with gRPC framework
  - GraphQL: Hasura Engine
  - Messaging: Apache Kafka with Go client
- **Databases**: PostgreSQL
- **Container Orchestration**: Docker Compose
- **Frontend**: HTML, CSS, JavaScript
- **API Technologies**:
  - REST
  - gRPC
  - GraphQL
  - Event-driven (Kafka)

---

## Future Improvements

- Implement authentication and authorization
- Add service discovery
- Implement circuit breakers for resilience
- Add comprehensive logging and monitoring
- Implement CI/CD pipeline
- Add unit and integration tests
