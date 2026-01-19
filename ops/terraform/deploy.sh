#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Azure Container Apps Deployment Script ===${NC}"

# Check prerequisites
command -v az >/dev/null 2>&1 || { echo -e "${RED}Error: Azure CLI is not installed${NC}" >&2; exit 1; }
command -v terraform >/dev/null 2>&1 || { echo -e "${RED}Error: Terraform is not installed${NC}" >&2; exit 1; }
command -v docker >/dev/null 2>&1 || { echo -e "${RED}Error: Docker is not installed${NC}" >&2; exit 1; }

# Check if logged in to Azure
echo -e "${YELLOW}Checking Azure login status...${NC}"
az account show >/dev/null 2>&1 || { echo -e "${RED}Error: Not logged in to Azure. Run 'az login' first${NC}" >&2; exit 1; }

# Check if terraform.tfvars exists
if [ ! -f "terraform.tfvars" ]; then
    echo -e "${YELLOW}terraform.tfvars not found. Creating from example...${NC}"
    cp terraform.tfvars.example terraform.tfvars
    echo -e "${RED}Please edit terraform.tfvars with your values, especially jwt_secret${NC}"
    exit 1
fi

# Initialize Terraform
echo -e "${YELLOW}Initializing Terraform...${NC}"
terraform init

# Plan
echo -e "${YELLOW}Planning Terraform deployment...${NC}"
terraform plan -out=tfplan

# Ask for confirmation
read -p "Do you want to apply this plan? (yes/no): " -r
if [[ ! $REPLY =~ ^[Yy]es$ ]]; then
    echo -e "${RED}Deployment cancelled${NC}"
    exit 1
fi

# Apply
echo -e "${YELLOW}Applying Terraform configuration...${NC}"
terraform apply tfplan

# Get outputs
echo -e "${YELLOW}Getting ACR details...${NC}"
ACR_LOGIN_SERVER=$(terraform output -raw acr_login_server)
ACR_USERNAME=$(terraform output -raw acr_admin_username)
ACR_PASSWORD=$(terraform output -raw acr_admin_password)
RESOURCE_GROUP=$(terraform output -raw resource_group_name)

# Login to ACR
echo -e "${YELLOW}Logging in to Azure Container Registry...${NC}"
echo "$ACR_PASSWORD" | docker login "$ACR_LOGIN_SERVER" -u "$ACR_USERNAME" --password-stdin

# Build and push image
echo -e "${YELLOW}Building Docker image...${NC}"
cd ../../app
docker build -t "$ACR_LOGIN_SERVER/go-tasks-app:latest" .

echo -e "${YELLOW}Pushing image to ACR...${NC}"
docker push "$ACR_LOGIN_SERVER/go-tasks-app:latest"

# Update container app
echo -e "${YELLOW}Updating Container App with new image...${NC}"
cd ../ops/terraform
az containerapp update \
  --name go-tasks-app \
  --resource-group "$RESOURCE_GROUP" \
  --image "$ACR_LOGIN_SERVER/go-tasks-app:latest"

# Get app URL
APP_URL=$(terraform output -raw container_app_url)

echo -e "${GREEN}=== Deployment Complete ===${NC}"
echo -e "${GREEN}Application URL: ${APP_URL}${NC}"
echo -e "${YELLOW}Note: It may take a few minutes for the app to be fully available${NC}"
