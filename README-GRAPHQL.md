# GraphQL Services Comparison

This project demonstrates two different approaches to implementing GraphQL services:

1. A custom Go-based GraphQL service 
2. A Hasura GraphQL service

Both services use the same PostgreSQL database and provide similar functionality, but with different implementation approaches.

## Running the Services

Start the services using Docker Compose:

```bash
docker-compose up -d
```

This will start:
- Custom GraphQL service on port 8083
- Hasura GraphQL service on port 8090
- API Gateway on port 8080 (which proxies to both services)

## Using the GraphQL Services

### Direct Access

- **Custom GraphQL Service**: http://localhost:8083/graphql
- **Hasura GraphQL Console**: http://localhost:8090/console

### Via API Gateway

- **Custom GraphQL Service**: http://localhost:8080/graphql
- **Hasura GraphQL Service**: http://localhost:8080/hasura

## Comparison UI

A comparison UI is provided to test both services side-by-side:

```
http://localhost:8080/compare-graphql.html
```

This UI allows you to:
- Query reviews by product ID on both services
- Create new reviews on both services
- Compare response formats and performance

## Schema Differences

### Custom GraphQL Service

```graphql
type Review {
  id: ID!
  productId: ID!
  userId: ID!
  username: String!
  rating: Int!
  comment: String!
  createdAt: String!
}

input CreateReviewInput {
  productId: ID!
  userId: ID!
  username: String!
  rating: Int!
  comment: String!
}
```

### Hasura GraphQL Service

```graphql
type reviews {
  id: uuid!
  product_id: String!
  user_id: String!
  username: String!
  rating: Int!
  comment: String!
  created_at: timestamptz!
  updated_at: timestamptz!
}

input reviews_insert_input {
  product_id: String!
  user_id: String!
  username: String!
  rating: Int!
  comment: String!
}
```

## Key Differences Between Approaches

| Feature | Custom GraphQL Service | Hasura |
|---------|------------------------|--------|
| Schema Definition | Manual (schema.graphql) | Auto-generated from DB |
| Resolvers | Manually coded | Auto-generated |
| Filters | Basic, custom-coded | Rich, auto-generated |
| Joins | Manual implementation | Automatic with relationships |
| Performance | Good for custom logic | Excellent for standard operations |
| Flexibility | Maximum code control | Less control, more productivity |
| Development Speed | Slower, more code | Faster, less code |
| Authentication | Custom implementation | Built-in, role-based | 