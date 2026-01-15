# Task Manager - Go + MongoDB + Templ

A full-featured task management web application with authentication, built with Go, MongoDB, and Templ. Features invite-only user registration and role-based access control.

## Features

- 🔐 **JWT + Session-based Authentication** - Secure authentication with HTTP-only cookies
- 👥 **Invite-only Registration** - Admins control who can join
- 📋 **Full CRUD for Tasks** - Create, read, update, and delete tasks
- 🎨 **Server-side Rendering** - Fast, modern UI with Templ
- 🔒 **Role-based Access Control** - Admin and user roles
- 🐳 **Docker Ready** - Complete Docker setup for local development
- ☁️ **Azure Container Apps** - Production-ready deployment configuration

## Prerequisites

- Docker and Docker Compose (recommended)
- Go 1.23+ (for local development)
- MongoDB 7+ (if running without Docker)

## Quick Start

### 1. Start Services

```bash
docker compose up --build
```

This will start:
- MongoDB on `localhost:27017`
- API server on `localhost:8080`

### 2. Create Admin User

```bash
# In a new terminal
docker exec -it tasks-api ./seed

# Or if running locally without Docker:
go run cmd/seed/main.go
```

Default admin credentials:
- **Email**: admin@example.com
- **Password**: admin123

To customize admin credentials, use environment variables:
```bash
ADMIN_EMAIL=your@email.com ADMIN_PASSWORD=yourpassword go run cmd/seed/main.go
```

### 3. Access the Application

Open http://localhost:8080 in your browser and log in with the admin credentials.

## Project Structure

```
.
├── cmd/
│   ├── server/          # Main application
│   └── seed/            # Admin user seeding tool
├── internal/
│   ├── auth/            # JWT, password hashing, middleware
│   ├── database/        # MongoDB repositories
│   ├── handlers/        # HTTP handlers
│   └── models/          # Data models
├── web/
│   ├── templates/       # Templ templates
│   └── static/          # CSS and static assets
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## Usage

### Admin Workflow

1. **Login** with admin credentials at `/login`
2. **Create Invites** at `/admin/invites`
3. **Copy Invite Link** and share with new users
4. **Manage Tasks** on the dashboard

### User Workflow

1. **Register** using an invite link
2. **View Tasks** on the dashboard
3. **Create Tasks** with title, description, status, and due date
4. **Edit/Delete** your own tasks

## API Endpoints

All API endpoints require authentication via JWT token in cookie or Authorization header.

### Tasks

```bash
# List all tasks (filtered by user)
GET /api/tasks

# Get specific task
GET /api/tasks/{id}

# Create task
POST /api/tasks
Content-Type: application/json
{
  "title": "Task title",
  "description": "Task description",
  "status": "pending",
  "due_date": "2026-01-20T00:00:00Z"
}

# Update task
PUT /api/tasks/{id}
Content-Type: application/json
{
  "title": "Updated title",
  "description": "Updated description",
  "status": "in_progress"
}

# Delete task
DELETE /api/tasks/{id}
```

### Authentication

```bash
# Login (returns JWT in cookie)
POST /login
Content-Type: application/x-www-form-urlencoded
email=user@example.com&password=password

# Logout
POST /logout

# Register (requires invite token)
POST /register/{token}
Content-Type: application/x-www-form-urlencoded
name=John Doe&email=john@example.com&password=password&confirm_password=password
```

## Environment Variables

Create a `.env` file or set environment variables:

```bash
# MongoDB
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=tasksdb

# Server
PORT=8080

# JWT Authentication
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRY=24h

# Admin Seed (optional)
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=admin123
ADMIN_NAME=Admin User
```

## Development

### Local Development (without Docker)

1. Start MongoDB:
```bash
docker run -d -p 27017:27017 mongo:7
```

2. Install templ:
```bash
go install github.com/a-h/templ/cmd/templ@latest
```

3. Generate templates:
```bash
templ generate
```

4. Run the server:
```bash
go run cmd/server/main.go
```

### Rebuilding Templates

When you modify `.templ` files:

```bash
templ generate
go run cmd/server/main.go
```

### Testing

```bash
# Build
go build ./cmd/server

# Run tests (if you add them)
go test ./...
```

## Deploying to Azure Container Apps

### 1. Build and Push Image

```bash
# Build for production
docker build -t your-registry.azurecr.io/task-manager:latest .

# Push to Azure Container Registry
docker push your-registry.azurecr.io/task-manager:latest
```

### 2. Create Container App

```bash
az containerapp create \
  --name task-manager \
  --resource-group your-rg \
  --environment your-env \
  --image your-registry.azurecr.io/task-manager:latest \
  --target-port 8080 \
  --ingress external \
  --env-vars \
    MONGODB_URI=your-mongodb-connection-string \
    MONGODB_DATABASE=tasksdb \
    JWT_SECRET=your-production-secret \
    JWT_EXPIRY=24h
```

### 3. Configure MongoDB

Use Azure Cosmos DB with MongoDB API or a managed MongoDB service:

```bash
# Example: Azure Cosmos DB connection string
MONGODB_URI=mongodb://account:key@account.mongo.cosmos.azure.com:10255/?ssl=true&replicaSet=globaldb
```

### 4. Seed Admin User

```bash
# SSH into container or use Azure Container Apps exec
az containerapp exec \
  --name task-manager \
  --resource-group your-rg \
  --command "./seed"
```

## Security Considerations

- ✅ Passwords hashed with bcrypt (cost 12)
- ✅ JWT stored in HTTP-only, SameSite=Strict cookies
- ✅ CSRF protection via SameSite cookies
- ✅ Single-use invite tokens with expiration
- ✅ User-scoped task access (users can only see their own tasks)
- ⚠️ Change `JWT_SECRET` in production
- ⚠️ Use HTTPS in production (set `Secure` flag on cookies)
- ⚠️ Consider rate limiting for login endpoint

## Task Status Values

- `pending` - Task not started
- `in_progress` - Task being worked on
- `completed` - Task finished

## Troubleshooting

### Docker Issues

```bash
# Clean rebuild
docker compose down -v
docker compose up --build

# View logs
docker compose logs -f api
```

### MongoDB Connection Issues

```bash
# Test MongoDB connection
docker exec -it mongodb mongosh
> show dbs
> use tasksdb
> db.users.find()
```

### Template Generation Issues

```bash
# Clean and regenerate
rm web/templates/*_templ.go
templ generate
```

## License

MIT

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request
