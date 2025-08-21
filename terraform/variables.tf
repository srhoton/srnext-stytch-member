variable "aws_region" {
  description = "AWS region for all resources"
  type        = string
  default     = "us-west-2"
}

variable "environment" {
  description = "Environment name (e.g., sandbox, dev, staging, prod)"
  type        = string
  default     = "sandbox"
}

variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "srnext-stytch-member"
}

# ALB Configuration
variable "alb_arn" {
  description = "ARN of the existing Application Load Balancer"
  type        = string
  default     = "arn:aws:elasticloadbalancing:us-west-2:345594586248:loadbalancer/app/external-private-alb/720e2b5474d3d602"
}

variable "alb_listener_port" {
  description = "ALB listener port to attach the target group to"
  type        = number
  default     = 443
}

variable "alb_path_patterns" {
  description = "Path patterns that should route to the Lambda"
  type        = list(string)
  default     = ["/members", "/members/*"]
}

variable "alb_rule_priority" {
  description = "Priority for the ALB listener rule"
  type        = number
  default     = 100
}

# DNS Configuration
variable "dns_name" {
  description = "DNS name for the service"
  type        = string
  default     = "srnext-stytch-member.sb.int.fullbayapi.com"
}

variable "dns_zone_name" {
  description = "Route53 hosted zone name"
  type        = string
  default     = "sb.int.fullbayapi.com"
}

# Lambda Configuration
variable "lambda_memory_size" {
  description = "Memory size for the Lambda function in MB"
  type        = number
  default     = 512

  validation {
    condition     = var.lambda_memory_size >= 128 && var.lambda_memory_size <= 10240
    error_message = "Lambda memory size must be between 128 and 10240 MB."
  }
}

variable "lambda_timeout" {
  description = "Timeout for the Lambda function in seconds"
  type        = number
  default     = 30

  validation {
    condition     = var.lambda_timeout >= 1 && var.lambda_timeout <= 900
    error_message = "Lambda timeout must be between 1 and 900 seconds."
  }
}

variable "lambda_reserved_concurrent_executions" {
  description = "Reserved concurrent executions for the Lambda function"
  type        = number
  default     = 100
}

variable "lambda_architecture" {
  description = "Lambda function architecture"
  type        = string
  default     = "arm64"

  validation {
    condition     = contains(["x86_64", "arm64"], var.lambda_architecture)
    error_message = "Lambda architecture must be either x86_64 or arm64."
  }
}

# Note: Stytch credentials are now retrieved from AWS Secrets Manager
# at srnext/stytchCredentials instead of being passed as variables

# CloudWatch Configuration
variable "log_retention_days" {
  description = "CloudWatch Logs retention in days"
  type        = number
  default     = 7

  validation {
    condition = contains([
      1, 3, 5, 7, 14, 30, 60, 90, 120, 150, 180, 365, 400, 545, 731, 1827, 3653
    ], var.log_retention_days)
    error_message = "Log retention days must be a valid CloudWatch Logs retention period."
  }
}

variable "enable_xray_tracing" {
  description = "Enable X-Ray tracing for the Lambda function"
  type        = bool
  default     = true
}

# Monitoring Configuration
variable "enable_cloudwatch_alarms" {
  description = "Enable CloudWatch alarms for the Lambda function"
  type        = bool
  default     = true
}

variable "alarm_email" {
  description = "Email address for CloudWatch alarm notifications"
  type        = string
  default     = ""
}

variable "error_rate_threshold" {
  description = "Error rate threshold for CloudWatch alarms (percentage)"
  type        = number
  default     = 5
}

variable "throttle_threshold" {
  description = "Throttle count threshold for CloudWatch alarms"
  type        = number
  default     = 10
}

# VPC Configuration (optional)
variable "vpc_subnet_ids" {
  description = "VPC subnet IDs for Lambda (leave empty for no VPC)"
  type        = list(string)
  default     = []
}

variable "vpc_security_group_ids" {
  description = "VPC security group IDs for Lambda (leave empty for no VPC)"
  type        = list(string)
  default     = []
}

# Tags
variable "default_tags" {
  description = "Default tags to apply to all resources"
  type        = map(string)
  default = {
    Environment = "sandbox"
    Project     = "srnext-stytch-member"
    ManagedBy   = "terraform"
    Owner       = "development-team"
  }
}