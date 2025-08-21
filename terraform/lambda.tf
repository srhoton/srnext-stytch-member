# Build the Lambda binary using a null resource
resource "null_resource" "lambda_build" {
  triggers = {
    always_run = timestamp()
  }

  provisioner "local-exec" {
    command     = "make clean build"
    working_dir = local.lambda_source_path
    environment = {
      GOOS   = "linux"
      GOARCH = var.lambda_architecture == "arm64" ? "arm64" : "amd64"
    }
  }
}

# Create the Lambda deployment package
data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = local.lambda_binary_path
  output_path = local.lambda_zip_path

  depends_on = [null_resource.lambda_build]
}

# Lambda function
resource "aws_lambda_function" "main" {
  function_name = local.lambda_function_name
  role          = aws_iam_role.lambda_execution.arn
  handler       = "bootstrap"
  runtime       = "provided.al2023"
  architectures = [var.lambda_architecture]

  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256

  memory_size                    = var.lambda_memory_size
  timeout                        = var.lambda_timeout
  reserved_concurrent_executions = var.lambda_reserved_concurrent_executions

  environment {
    variables = {
      # Stytch configuration from Secrets Manager
      STYTCH_PROJECT_ID = local.stytch_credentials.project_id
      STYTCH_SECRET     = local.stytch_credentials.secret
      STYTCH_ENV        = "test" # The credentials are for test environment

      # Configuration
      LOG_LEVEL       = "info"
      REQUEST_TIMEOUT = "30"
      MAX_RETRIES     = "3"

      # Additional configuration
      ENVIRONMENT = var.environment
    }
  }

  # VPC configuration (optional)
  dynamic "vpc_config" {
    for_each = length(var.vpc_subnet_ids) > 0 ? [1] : []
    content {
      subnet_ids         = var.vpc_subnet_ids
      security_group_ids = var.vpc_security_group_ids
    }
  }

  # Enable X-Ray tracing
  dynamic "tracing_config" {
    for_each = var.enable_xray_tracing ? [1] : []
    content {
      mode = "Active"
    }
  }

  # Dead letter queue (optional)
  # dead_letter_config {
  #   target_arn = aws_sqs_queue.dlq.arn
  # }

  tags = merge(
    local.common_tags,
    {
      Name = local.lambda_function_name
    }
  )

  depends_on = [
    aws_iam_role_policy_attachment.lambda_basic_execution,
    aws_iam_role_policy_attachment.lambda_cloudwatch_logs,
    aws_cloudwatch_log_group.lambda,
    null_resource.lambda_build
  ]
}

# Lambda permission for ALB to invoke
resource "aws_lambda_permission" "alb_invoke" {
  statement_id  = "AllowExecutionFromALB"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.main.function_name
  principal     = "elasticloadbalancing.amazonaws.com"
  source_arn    = aws_lb_target_group.lambda.arn
}

# Lambda alias for production-like deployment
resource "aws_lambda_alias" "main" {
  name             = var.environment
  description      = "Alias for ${var.environment} environment"
  function_name    = aws_lambda_function.main.function_name
  function_version = aws_lambda_function.main.version

  # Optionally configure weighted routing for canary deployments
  # routing_config {
  #   additional_version_weights = {
  #     "2" = 0.1 # 10% traffic to version 2
  #   }
  # }
}