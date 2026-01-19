variable "resource_group_name" {
  description = "Name of the resource group"
  type        = string
  default     = "rg-go-app"
}

variable "location" {
  description = "Azure region for resources"
  type        = string
  default     = "East US"
}

variable "cosmos_account_name" {
  description = "Name of the CosmosDB account"
  type        = string
  default     = "cosmos-go-app"
}

variable "database_name" {
  description = "Name of the MongoDB database"
  type        = string
  default     = "tasksdb"
}

variable "acr_name" {
  description = "Name of the Azure Container Registry"
  type        = string
  default     = "acrgoapp"
}

variable "log_analytics_name" {
  description = "Name of the Log Analytics workspace"
  type        = string
  default     = "log-go-app"
}

variable "key_vault_name" {
  description = "Name of the Key Vault"
  type        = string
  default     = "kv-go-app"
}

variable "container_app_environment_name" {
  description = "Name of the Container Apps environment"
  type        = string
  default     = "cae-go-app"
}

variable "app_name" {
  description = "Name of the container app"
  type        = string
  default     = "go-tasks-app"
}

variable "image_tag" {
  description = "Docker image tag"
  type        = string
  default     = "latest"
}

variable "jwt_secret" {
  description = "JWT secret for authentication"
  type        = string
  sensitive   = true
  default     = ""
}

variable "jwt_expiry" {
  description = "JWT token expiry duration"
  type        = string
  default     = "24h"
}

variable "min_replicas" {
  description = "Minimum number of container replicas"
  type        = number
  default     = 1
}

variable "max_replicas" {
  description = "Maximum number of container replicas"
  type        = number
  default     = 3
}

variable "container_cpu" {
  description = "CPU allocation for container"
  type        = number
  default     = 0.5
}

variable "container_memory" {
  description = "Memory allocation for container"
  type        = string
  default     = "1Gi"
}

variable "tags" {
  description = "Tags to apply to resources"
  type        = map(string)
  default = {
    Environment = "Production"
    ManagedBy   = "Terraform"
  }
}
