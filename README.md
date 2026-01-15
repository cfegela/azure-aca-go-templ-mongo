# Tasks CRUD API

A RESTful API for managing tasks built with Go and MongoDB, containerized for local development and Azure Container Apps deployment.

## Features

- Full CRUD operations for tasks
- MongoDB for data persistence
- Docker and docker-compose for local development
- Health check endpoint for container orchestration
- CORS enabled
- Graceful shutdown handling

## Prerequisites

- Docker and Docker Compose
- Go 1.22+ (for local development without Docker)

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Start the API and MongoDB
docker-compose up --build

# The API will be available at http://localhost:8080
```

### Using Go Directly

```bash
# Install dependencies
go mod download

# Start MongoDB separately (required)
docker run -d -p 27017:27017 mongo:7

# Run the API
go run cmd/server/main.go
```

## API Endpoints

### Create a Task
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete project",
    "description": "Finish the Go API implementation",
    "status": "pending",
    "due_date": "2026-01-20T00:00:00Z"
  }'
```

### List All Tasks
```bash
curl http://localhost:8080/api/tasks
```

### Get a Specific Task
```bash
curl http://localhost:8080/api/tasks/{id}
```

### Update a Task
```bash
curl -X PUT http://localhost:8080/api/tasks/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated title",
    "description": "Updated description",
    "status": "in_progress"
  }'
```

### Delete a Task
```bash
curl -X DELETE http://localhost:8080/api/tasks/{id}
```

### Health Check
```bash
curl http://localhost:8080/health
```

## Task Model

```json
{
  "id": "507f1f77bcf86cd799439011",
  "title": "Task title",
  "description": "Task description",
  "status": "pending",
  "due_date": "2026-01-20T00:00:00Z",
  "created_at": "2026-01-15T10:00:00Z",
  "updated_at": "2026-01-15T10:00:00Z"
}
```

### Status Values
- `pending` - Task is not started
- `in_progress` - Task is being worked on
- `completed` - Task is finished

## Environment Variables

Create a `.env` file from `.env.example`:

```bash
cp .env.example .env
```

Available variables:
- `MONGODB_URI` - MongoDB connection string (default: `mongodb://localhost:27017`)
- `MONGODB_DATABASE` - Database name (default: `tasksdb`)
- `PORT` - Server port (default: `8080`)

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── handlers/
│   │   └── tasks.go             # HTTP handlers
│   ├── models/
│   │   └── task.go              # Task model
│   └── database/
│       └── mongodb.go           # MongoDB operations
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## Deploying to Azure Container Apps

1. Build and push the Docker image:
```bash
docker build -t your-registry.azurecr.io/tasks-api:latest .
docker push your-registry.azurecr.io/tasks-api:latest
```

2. Create Azure Container App with MongoDB connection:
```bash
az containerapp create \
  --name tasks-api \
  --resource-group your-rg \
  --environment your-env \
  --image your-registry.azurecr.io/tasks-api:latest \
  --target-port 8080 \
  --ingress external \
  --env-vars MONGODB_URI=your-mongodb-uri MONGODB_DATABASE=tasksdb
```

3. Configure MongoDB (use Azure Cosmos DB with MongoDB API or a managed MongoDB service)

## Development

### Run Tests
```bash
go test ./...
```

### Build Locally
```bash
go build -o bin/server ./cmd/server
./bin/server
```

## License

MIT
