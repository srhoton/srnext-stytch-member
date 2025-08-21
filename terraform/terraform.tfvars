# Terraform Variables for Stytch Member Management Lambda
# Sandbox Environment Configuration

# Project Configuration
project_name = "srnext-stytch-member"
environment  = "sandbox"
aws_region   = "us-west-2"

# Stytch Configuration
# Note: Stytch credentials are now retrieved automatically from AWS Secrets Manager
# Secret name: srnext/stytchCredentials

# ALB Configuration
alb_arn           = "arn:aws:elasticloadbalancing:us-west-2:345594586248:loadbalancer/app/external-private-alb/720e2b5474d3d602"
alb_listener_port = 443

# DNS Configuration
dns_zone_name   = "sb.int.fullbayapi.com"
dns_record_name = "srnext-stytch-member.sb.int.fullbayapi.com"

# Lambda Configuration
lambda_memory_size = 256
lambda_timeout     = 30

# CloudWatch Configuration
log_retention_days       = 7
enable_cloudwatch_alarms = true
alarm_email              = ""

# X-Ray Tracing
enable_xray_tracing = false

# VPC Configuration (leave empty for no VPC)
vpc_subnet_ids      = []
vpc_security_groups = []

# Tags
default_tags = {
  Owner       = "Steve Rhoton"
  ManagedBy   = "terraform"
  Environment = "sandbox"
  Project     = "srnext-stytch-member"
}