# Postman Collection for Go Microservices Project

This Postman collection provides a complete set of requests to test all APIs in the Go Microservices Project. The collection includes tests for the Products API and health checks.

## How to Import the Collection

1. Make sure you have [Postman](https://www.postman.com/downloads/) installed
2. Open Postman
3. Click on "Import" button in the top left corner
4. Select "File" and choose the `postman-collection.json` file
5. Click "Import"

## Setting Up Environment Variables

Before using the collection, you should set up an environment with the following variables:

1. Click on "Environments" in the sidebar
2. Click "+" to create a new environment
3. Name it "Go Microservices Local"
4. Add the following variables:
   - `baseUrl`: `http://localhost:8080`
   - `product_id`: Leave empty (will be populated after running "List All Products")
   - `new_product_id`: Leave empty (will be populated after running "Create New Product")
5. Click "Save"
6. Select the "Go Microservices Local" environment from the dropdown in the top right corner

## Using the Collection

The collection is organized into folders by service. Here's how to use it:

### Health Checks

1. Start with the "API Gateway Health" request to ensure the API Gateway is running

### Products API

Execute the requests in this order for the best experience:

1. **List All Products**: Gets all products and automatically stores the ID of the first product in the `product_id` variable
2. **Get Product by ID**: Retrieves the product using the stored ID
3. **Create New Product**: Creates a new product and stores its ID in the `new_product_id` variable
4. **Update Product**: Updates the newly created product
5. **Delete Product**: Deletes the product you created
6. **Filter Products by Category**: Lists products filtered by the "Test" category

## Test Scripts

Each request includes test scripts that:
- Validate the response status code
- Check that the response body contains the expected fields
- Automatically set environment variables for use in subsequent requests

## Troubleshooting

If requests fail, check that:
1. All microservices are running (use `./start-project.sh`)
2. The ports are correctly configured (8080 for API Gateway)
3. You have selected the correct environment in Postman

## Extending the Collection

To add more requests to the collection:
1. Duplicate an existing request
2. Update the method, URL, and body as needed
3. Modify the test script to check for the expected response
4. Save your changes

## API Reference

### Products API

- `GET /api/products/`: List all products
- `GET /api/products/{id}`: Get a product by ID
- `POST /api/products/`: Create a new product
- `PUT /api/products/{id}`: Update a product
- `DELETE /api/products/{id}`: Delete a product
- `GET /api/products/?category={category}`: Filter products by category

### Health Checks

- `GET /health`: API Gateway health check 