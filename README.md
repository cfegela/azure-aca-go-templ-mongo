# Task Manager - Go + MongoDB + Templ

A full-featured task management web application with authentication, built with Go, MongoDB, and Templ. Features invite-only user registration and role-based access control.

## Repository Structure

```
.
├── app/                    # Application code
│   ├── cmd/
│   │   ├── server/        # Main application
│   │   └── seed/          # Admin user seeding tool
│   ├── internal/
│   │   ├── auth/          # JWT, password hashing, middleware
│   │   ├── database/      # MongoDB repositories
│   │   ├── handlers/      # HTTP handlers
│   │   └── models/        # Data models
│   ├── web/
│   │   ├── templates/     # Templ templates
│   │   └── static/        # CSS and static assets
│   ├── Dockerfile
│   ├── docker-compose.yml
│   └── .env.example
├── ops/
│   └── terraform/         # Infrastructure as Code
│       ├── main.tf
│       ├── variables.tf
│       ├── outputs.tf
│       ├── deploy.sh      # Automated deployment script
│       └── terraform.tfvars.example
└── README.md
```

## Features

- 🔐 **JWT + Session-based Authentication** - Secure authentication with HTTP-only cookies
- 👥 **Invite-only Registration** - Admins control who can join
- 📋 **Full CRUD for Tasks** - Create, read, update, and delete tasks
- 🎨 **Server-side Rendering** - Fast, modern UI with Templ
- 🔒 **Role-based Access Control** - Admin and user roles
- 🐳 **Docker Ready** - Complete Docker setup for local development
- ☁️ **Azure Container Apps** - Production-ready deployment configuration
- 🏥 **Health Checks** - Built-in health check endpoints for monitoring
- 🔄 **Auto-Seeding** - Automatic admin user creation on container startup

## Prerequisites

- Docker and Docker Compose (recommended)
- Go 1.23+ (for local development)
- MongoDB 7+ (if running without Docker)

## Quick Start

### 1. Start Services

```bash
cd app
docker compose up --build
```

This will start:
- MongoDB on `localhost:27017`
- API server on `localhost:8080`

The container automatically runs the seed script on startup, creating an admin user if one doesn't exist.

### 2. Default Admin Credentials

The application automatically creates an admin user on first startup:
- **Email**: admin@example.com
- **Password**: admin123

To customize admin credentials, set environment variables before starting:
```bash
ADMIN_EMAIL=your@email.com ADMIN_PASSWORD=yourpassword docker compose up --build
```

Or manually run the seed tool:
```bash
# Inside the container
docker exec -it tasks-api ./seed

# Or locally without Docker
cd app
go run cmd/seed/main.go
```

### 3. Access the Application

Open http://localhost:8080 in your browser and log in with the admin credentials.

### 4. Health Check

The application includes a health check endpoint:
```bash
curl http://localhost:8080/health
```

## Application Architecture

### Core Components

- **Server** (`app/cmd/server/main.go`): Main HTTP server with routing and middleware
- **Seed Tool** (`app/cmd/seed/main.go`): Admin user initialization utility
- **Authentication** (`app/internal/auth/`): JWT generation, password hashing, middleware
- **Database** (`app/internal/database/`): MongoDB repositories for Users, Tasks, and Invites
- **Handlers** (`app/internal/handlers/`): HTTP request handlers for API and pages
- **Models** (`app/internal/models/`): Data structures for User, Task, and Invite
- **Templates** (`app/web/templates/`): Templ-based server-side rendered UI

### Server Features

- **Graceful Shutdown**: Handles SIGINT with 10-second timeout
- **CORS Support**: Configured for cross-origin requests
- **Request Timeouts**: 15s read/write, 60s idle
- **Database Indexes**: Automatic index creation for performance
- **Health Endpoint**: `/health` returns JSON status
- **Static File Serving**: `/static/` for CSS and assets

### Routes

#### Public Routes
- `GET /login` - Login page
- `POST /login` - Login form submission
- `POST /logout` - Logout
- `GET /register/{token}` - Registration page with invite token
- `POST /register/{token}` - Registration form submission
- `GET /health` - Health check endpoint

#### Protected Routes (Require Authentication)
- `GET /` - Dashboard with task list
- `GET /tasks/new` - New task form
- `POST /tasks` - Create task
- `GET /tasks/{id}/edit` - Edit task form
- `POST /tasks/{id}` - Update task
- `POST /tasks/{id}/delete` - Delete task

#### Admin Routes (Require Admin Role)
- `GET /admin/invites` - Invite management page
- `POST /admin/invites` - Create new invite

#### API Routes (Require Authentication)
- `GET /api/tasks` - List all user's tasks (JSON)
- `GET /api/tasks/{id}` - Get specific task (JSON)
- `POST /api/tasks` - Create task (JSON)
- `PUT /api/tasks/{id}` - Update task (JSON)
- `DELETE /api/tasks/{id}` - Delete task (JSON)

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

Create a `.env` file in the `app/` directory or set environment variables:

```bash
# MongoDB
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=tasksdb

# Server
PORT=8080

# JWT Authentication
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRY=24h

# Admin Seed (optional - used by seed tool)
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=admin123
ADMIN_NAME=Admin User
```

See `app/.env.example` for a template.

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

3. Navigate to the app directory:
```bash
cd app
```

4. Generate templates:
```bash
templ generate
```

5. Run the seed tool (first time only):
```bash
go run cmd/seed/main.go
```

6. Run the server:
```bash
go run cmd/server/main.go
```

### Rebuilding Templates

When you modify `.templ` files:

```bash
cd app
templ generate
go run cmd/server/main.go
```

### Testing

```bash
cd app

# Build server
go build ./cmd/server

# Build seed tool
go build ./cmd/seed

# Run tests (if you add them)
go test ./...
```

## Deploying to Azure Container Apps

### Option 1: Automated Deployment with deploy.sh

The easiest way to deploy is using the automated deployment script:

```bash
cd ops/terraform

# First time: Copy and configure variables
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your values (especially jwt_secret)

# Run the automated deployment
./deploy.sh
```

The `deploy.sh` script will:
1. Check prerequisites (Azure CLI, Terraform, Docker)
2. Verify Azure login status
3. Initialize and plan Terraform deployment
4. Apply infrastructure changes
5. Build and push Docker image to ACR
6. Update the Container App with the new image
7. Display the application URL

### Option 2: Manual Deployment

#### 1. Build and Push Image

```bash
# Build for production
cd app
docker build -t your-registry.azurecr.io/go-tasks-app:latest .

# Push to Azure Container Registry
docker push your-registry.azurecr.io/go-tasks-app:latest
```

#### 2. Create Container App

```bash
az containerapp create \
  --name go-tasks-app \
  --resource-group your-rg \
  --environment your-env \
  --image your-registry.azurecr.io/go-tasks-app:latest \
  --target-port 8080 \
  --ingress external \
  --env-vars \
    MONGODB_URI=your-mongodb-connection-string \
    MONGODB_DATABASE=tasksdb \
    JWT_SECRET=your-production-secret \
    JWT_EXPIRY=24h
```

#### 3. Configure MongoDB

Use Azure Cosmos DB with MongoDB API (automatically configured by Terraform):

```bash
# Example: Azure Cosmos DB connection string
MONGODB_URI=mongodb://account:key@account.mongo.cosmos.azure.com:10255/?ssl=true&replicaSet=globaldb
```

#### 4. Verify Deployment

The container automatically seeds the admin user on startup. Check the application health:

```bash
curl https://your-app.azurecontainerapps.io/health
```

## Docker Configuration

### Multi-stage Build

The `app/Dockerfile` uses a multi-stage build for optimal image size:

1. **Builder Stage** (golang:1-alpine)
   - Installs dependencies
   - Generates Templ templates
   - Builds both `server` and `seed` binaries

2. **Runtime Stage** (alpine:latest)
   - Minimal base image with ca-certificates
   - Copies compiled binaries and static files
   - Runs startup script that seeds admin user then starts server

### Automatic Admin Seeding

The container includes a startup script (`/root/start.sh`) that:
1. Runs the seed tool to create admin user (skips if already exists)
2. Starts the main application server

This ensures the application is ready to use immediately after deployment without manual intervention.

### Docker Compose

The `app/docker-compose.yml` includes:
- **MongoDB 7** with persistent volume
- **API server** with health check dependency
- **Automatic networking** between services
- **Environment variables** for configuration

Health checks ensure the API doesn't start until MongoDB is ready.

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
cd app

# Clean rebuild
docker compose down -v
docker compose up --build

# View logs
docker compose logs -f api

# Check if seed script ran successfully
docker logs tasks-api | grep -i seed
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
cd app

# Clean and regenerate
rm web/templates/*_templ.go
templ generate
```

### Health Check Issues

```bash
# Test health endpoint locally
curl http://localhost:8080/health

# Test health endpoint on Azure
curl https://your-app.azurecontainerapps.io/health
```

## Infrastructure & Terraform Deployment

This project includes Terraform configuration for deploying to Azure Container Apps with CosmosDB (MongoDB API).

### Architecture

The infrastructure includes:

- **Resource Group**: Container for all Azure resources
- **CosmosDB Account**: Serverless MongoDB API database
- **Azure Container Registry (ACR)**: Private container registry
- **Container Apps Environment**: Runtime environment for containers
- **Container App**: The deployed application
- **Key Vault**: Secure storage for secrets (JWT secret)
- **Log Analytics Workspace**: Centralized logging and monitoring
- **Managed Identity**: Secure authentication between Azure services

### Prerequisites for Infrastructure Deployment

1. **Azure CLI**: [Install Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)
2. **Terraform**: [Install Terraform](https://www.terraform.io/downloads.html) (>= 1.0)
3. **Docker**: For building and pushing container images
4. **Azure Subscription**: Active Azure subscription

### Terraform Setup Instructions

#### Quick Start with deploy.sh

```bash
cd ops/terraform

# Copy and edit configuration
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your values

# Run automated deployment
./deploy.sh
```

#### Manual Terraform Deployment

##### 1. Authenticate with Azure

```bash
az login
az account set --subscription <your-subscription-id>
```

##### 2. Configure Variables

Copy the example tfvars file and customize it:

```bash
cd ops/terraform
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars` with your values:
- Update `jwt_secret` with a secure random string
- Modify `location` if needed
- Adjust resource names and scaling parameters
- Configure `tags` for resource management

##### 3. Initialize Terraform

```bash
terraform init
```

##### 4. Plan Deployment

```bash
terraform plan
```

Review the planned changes to ensure everything looks correct.

##### 5. Apply Configuration

```bash
terraform apply
```

Type `yes` when prompted to create the resources.

##### 6. Build and Push Docker Image

After the infrastructure is created, get the ACR credentials:

```bash
# Get ACR login server
ACR_LOGIN_SERVER=$(terraform output -raw acr_login_server)

# Get ACR credentials
ACR_USERNAME=$(terraform output -raw acr_admin_username)
ACR_PASSWORD=$(terraform output -raw acr_admin_password)

# Login to ACR
echo $ACR_PASSWORD | docker login $ACR_LOGIN_SERVER -u $ACR_USERNAME --password-stdin

# Build and push the image
cd ../../app
docker build -t $ACR_LOGIN_SERVER/go-tasks-app:latest .
docker push $ACR_LOGIN_SERVER/go-tasks-app:latest
```

##### 7. Update Container App

After pushing the image, update the container app to use the new image:

```bash
cd ../ops/terraform
terraform apply -auto-approve
```

Or use Azure CLI to create a new revision:

```bash
az containerapp update \
  --name go-tasks-app \
  --resource-group $(terraform output -raw resource_group_name) \
  --image $ACR_LOGIN_SERVER/go-tasks-app:latest
```

##### 8. Access Your Application

Get the application URL:

```bash
terraform output container_app_url
```

Visit the URL in your browser to access your application.

### Terraform Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `resource_group_name` | Name of the resource group | `rg-go-app` |
| `location` | Azure region | `East US` |
| `cosmos_account_name` | CosmosDB account name | `cosmos-go-app` |
| `database_name` | MongoDB database name | `tasksdb` |
| `acr_name` | Azure Container Registry name | `acrgoapp` |
| `log_analytics_name` | Log Analytics workspace name | `log-go-app` |
| `key_vault_name` | Key Vault name | `kv-go-app` |
| `container_app_environment_name` | Container Apps environment name | `cae-go-app` |
| `app_name` | Container app name | `go-tasks-app` |
| `image_tag` | Docker image tag | `latest` |
| `jwt_secret` | JWT secret (sensitive) | `` |
| `jwt_expiry` | JWT token expiry | `24h` |
| `min_replicas` | Minimum replicas | `1` |
| `max_replicas` | Maximum replicas | `3` |
| `container_cpu` | CPU allocation | `0.5` |
| `container_memory` | Memory allocation | `1Gi` |
| `tags` | Resource tags | `{Environment = "Production", ManagedBy = "Terraform"}` |

### Terraform Outputs

| Output | Description |
|--------|-------------|
| `resource_group_name` | Name of the created resource group |
| `container_app_url` | URL of the deployed application |
| `container_app_fqdn` | FQDN of the container app |
| `acr_login_server` | ACR login server address |
| `acr_admin_username` | ACR admin username (sensitive) |
| `acr_admin_password` | ACR admin password (sensitive) |
| `cosmos_endpoint` | CosmosDB endpoint |
| `cosmos_connection_string` | CosmosDB connection string (sensitive) |
| `key_vault_name` | Key Vault name |
| `key_vault_uri` | Key Vault URI |
| `log_analytics_workspace_id` | Log Analytics workspace ID |
| `managed_identity_client_id` | Managed identity client ID |
| `managed_identity_principal_id` | Managed identity principal ID |

### Monitoring Infrastructure

#### View Logs

```bash
az containerapp logs show \
  --name go-tasks-app \
  --resource-group $(terraform output -raw resource_group_name) \
  --follow
```

#### View Metrics

Access Log Analytics workspace in Azure Portal or use:

```bash
az monitor log-analytics workspace show \
  --workspace-name $(terraform output -raw log_analytics_workspace_id | xargs basename) \
  --resource-group $(terraform output -raw resource_group_name)
```

### Scaling Configuration

The application is configured with:
- Min replicas: 1 (configurable via `min_replicas`)
- Max replicas: 3 (configurable via `max_replicas`)
- Auto-scaling based on HTTP requests and CPU usage

To modify scaling:

```bash
az containerapp update \
  --name go-tasks-app \
  --resource-group $(terraform output -raw resource_group_name) \
  --min-replicas 2 \
  --max-replicas 5
```

### Cost Optimization

- **CosmosDB**: Using serverless mode (pay-per-request)
- **Container Apps**: Pay only for running containers
- **ACR**: Basic tier for small-scale usage
- **Log Analytics**: 30-day retention

To reduce costs:
- Decrease `max_replicas` if traffic is low
- Reduce `container_cpu` and `container_memory` if possible
- Consider disabling auto-scaling for dev environments

### Infrastructure Cleanup

To destroy all resources:

```bash
terraform destroy
```

Type `yes` when prompted.

### Infrastructure Troubleshooting

#### Container App Not Starting

1. Check logs:
   ```bash
   az containerapp logs show --name go-tasks-app --resource-group <rg-name> --follow
   ```

2. Verify image exists in ACR:
   ```bash
   az acr repository show-tags --name <acr-name> --repository go-tasks-app
   ```

#### CosmosDB Connection Issues

1. Verify connection string:
   ```bash
   terraform output cosmos_connection_string
   ```

2. Check if MongoDB API is enabled in CosmosDB account

#### Key Vault Access Issues

Ensure the managed identity has proper permissions:
```bash
az role assignment list --assignee <managed-identity-principal-id> --scope <key-vault-id>
```

### Security Best Practices (Infrastructure)

1. **JWT Secret**: Always use a strong, randomly generated secret
2. **RBAC**: Key Vault uses RBAC instead of access policies
3. **Managed Identity**: No credentials stored in code
4. **Secrets**: Sensitive outputs marked as sensitive
5. **Network**: Consider adding virtual network integration for production

### CI/CD Integration

#### Using the Deploy Script

The repository includes `ops/terraform/deploy.sh` which automates the entire deployment process:

```bash
cd ops/terraform
./deploy.sh
```

The script handles:
1. Prerequisites validation (Azure CLI, Terraform, Docker)
2. Azure authentication check
3. Terraform initialization and planning
4. Infrastructure provisioning
5. Docker image build and push to ACR
6. Container App revision update
7. Health verification and URL display

#### Manual CI/CD Workflow

For custom CI/CD pipelines (GitHub Actions, Azure DevOps, etc.), follow this workflow:

1. Authenticate with Azure
2. Build Docker image from `app/` directory
3. Push to ACR
4. Run `terraform apply` or use `az containerapp update`
5. Container Apps automatically creates a new revision
6. Verify health endpoint: `/health`

#### Example GitHub Actions Workflow

```yaml
name: Deploy to Azure

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Azure Login
        uses: azure/login@v1
        with:
          creds: ${{ secrets.AZURE_CREDENTIALS }}

      - name: Build and Deploy
        run: |
          cd ops/terraform
          ./deploy.sh
```

## License

MIT

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request
