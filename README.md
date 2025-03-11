# Retail Pulse Image Processor

## Description

The Retail Pulse Image Processor is a Go microservice designed to process thousands of images collected from retail stores. The service receives jobs containing image URLs and store IDs, downloads images, computes their perimeters, and stores the results.

Key features:
- Asynchronous job processing
- Concurrent image processing
- REST API for job submission and status checking
- Store validation against a master store list
- In-memory data storage for simplicity

## Assumptions

1. The `store-master.json` file is available in the root directory of the project
2. Images are publicly accessible via the provided URLs
3. Images are in standard formats (JPEG, PNG, GIF)
4. In-memory storage is sufficient for the assignment (in a production environment, a database would be used)
5. Jobs and results do not need to persist after service restart
6. The service is designed to run on a single instance (not distributed)

## Setup Instructions

### Prerequisites

- Go 1.19 or later
- Docker and Docker Compose (optional)

### Installation

#### Using Go directly

1. Clone the repository
2. Navigate to the project directory
3. Run the application:

```bash
go mod download
go run ./cmd/main.go
```

#### Using Docker

1. Clone the repository
2. Navigate to the project directory
3. Build and run with Docker Compose:

```bash
docker-compose up --build
```

The service will be available at http://localhost:8080

## API Usage

### Submit a Job

```bash
curl -X POST http://localhost:8080/api/submit/ \
  -H "Content-Type: application/json" \
  -d '{
    "count": 2,
    "visits": [
      {
        "store_id": "S00339218",
        "image_url": [
          "https://www.gstatic.com/webp/gallery/2.jpg",
          "https://www.gstatic.com/webp/gallery/3.jpg"
        ],
        "visit_time": "2023-06-15T10:30:00Z"
      },
      {
        "store_id": "S01408764",
        "image_url": [
          "https://www.gstatic.com/webp/gallery/3.jpg"
        ],
        "visit_time": "2023-06-15T11:15:00Z"
      }
    ]
  }'
```

### Check Job Status

```bash
curl http://localhost:8080/api/status?jobid=<JOB_ID>
```

## Testing

Run the Go tests:

```bash
go test ./...
```

## Work Environment

- **Computer/OS**: MacBook Pro, macOS Ventura 13.4
- **IDE**: Visual Studio Code with Go extension
- **Libraries**: Standard Go libraries only
- **Tools**: Docker, Docker Compose, Go 1.19

## Future Improvements

With more time, I would implement the following improvements:

1. **Persistent Storage**: Replace in-memory storage with a database like PostgreSQL or MongoDB
2. **Distributed Processing**: Add support for distributed processing using a message queue (RabbitMQ, Kafka)
3. **Enhanced Error Handling**: More detailed error handling and recovery mechanisms
4. **Metrics & Monitoring**: Add Prometheus metrics and logging with ELK stack
5. **Caching**: Add caching for frequently accessed data to improve performance
6. **Rate Limiting**: Implement rate limiting for the API endpoints
7. **Authentication & Authorization**: Add API key or OAuth-based authentication
8. **Swagger Documentation**: Add OpenAPI/Swagger documentation for the API
9. **Comprehensive Testing**: Add unit tests, integration tests, and load tests
10. **Image Processing Optimization**: Optimize image processing with parallel downloading and processing
11. **Health Checks**: Add health check endpoints for monitoring
12. **Pagination**: Add pagination for large result sets
13. **Circuit Breaker**: Implement circuit breaker patterns for external service calls