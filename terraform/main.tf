# Main Terraform Configuration for Stytch Member Management Lambda
# This file serves as the primary entrypoint for the module

# Module Information
# ==================
# This Terraform module deploys a serverless member management system
# for Stytch B2B organizations using AWS Lambda, DynamoDB, and ALB.
#
# The infrastructure is organized across the following files:
# - versions.tf: Terraform and provider requirements
# - variables.tf: Input variable definitions
# - locals.tf: Local values and computed variables
# - data.tf: Data sources for existing resources
# - lambda.tf: Lambda function and related resources
# - dynamodb.tf: DynamoDB table configuration
# - iam.tf: IAM roles and policies
# - alb.tf: ALB target group and listener rules
# - route53.tf: DNS configuration
# - cloudwatch.tf: Monitoring and alarms
# - outputs.tf: Output values

# Module Usage
# ============
# 1. Copy terraform.tfvars.example to terraform.tfvars
# 2. Configure required variables (Stytch credentials, etc.)
# 3. Run: terraform init
# 4. Run: terraform plan
# 5. Run: terraform apply

# Environment
# ===========
# This configuration is designed for the sandbox environment
# and integrates with existing infrastructure (ALB, Route53 zone).