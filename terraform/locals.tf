locals {
  # Naming conventions
  name_prefix = "${var.project_name}-${var.environment}"

  # Lambda function name
  lambda_function_name = "${local.name_prefix}-lambda"

  # CloudWatch log group name
  log_group_name = "/aws/lambda/${local.lambda_function_name}"

  # Target group name
  target_group_name = "${local.name_prefix}-tg"

  # Lambda source path
  lambda_source_path = "${path.module}/../lambda"

  # Lambda binary path
  lambda_binary_path = "${local.lambda_source_path}/build/bootstrap"

  # Lambda deployment package
  lambda_zip_path = "${path.module}/lambda-deployment.zip"

  # Parse Stytch credentials from Secrets Manager
  stytch_credentials = jsondecode(data.aws_secretsmanager_secret_version.stytch_credentials.secret_string)

  # Common tags
  common_tags = merge(
    var.default_tags,
    {
      Name        = local.name_prefix
      Environment = var.environment
      Project     = var.project_name
      ManagedBy   = "terraform"
    }
  )

}