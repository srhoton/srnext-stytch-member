output "lambda_function_name" {
  description = "Name of the Lambda function"
  value       = aws_lambda_function.main.function_name
}

output "lambda_function_arn" {
  description = "ARN of the Lambda function"
  value       = aws_lambda_function.main.arn
}

output "lambda_function_version" {
  description = "Version of the Lambda function"
  value       = aws_lambda_function.main.version
}

output "lambda_function_alias_arn" {
  description = "ARN of the Lambda function alias"
  value       = aws_lambda_alias.main.arn
}

output "lambda_execution_role_arn" {
  description = "ARN of the Lambda execution role"
  value       = aws_iam_role.lambda_execution.arn
}

output "target_group_arn" {
  description = "ARN of the ALB target group"
  value       = aws_lb_target_group.lambda.arn
}

output "cloudwatch_log_group_name" {
  description = "Name of the CloudWatch log group"
  value       = aws_cloudwatch_log_group.lambda.name
}

output "cloudwatch_dashboard_url" {
  description = "URL of the CloudWatch dashboard"
  value       = "https://${var.aws_region}.console.aws.amazon.com/cloudwatch/home?region=${var.aws_region}#dashboards:name=${aws_cloudwatch_dashboard.main.dashboard_name}"
}

output "service_endpoint" {
  description = "HTTPS endpoint for the service"
  value       = "https://${var.dns_name}"
}

output "alb_dns_name" {
  description = "DNS name of the ALB"
  value       = data.aws_lb.existing.dns_name
}

output "route53_record_name" {
  description = "Route53 record name"
  value       = aws_route53_record.main.name
}

output "route53_record_fqdn" {
  description = "Route53 record FQDN"
  value       = aws_route53_record.main.fqdn
}

output "sns_topic_arn" {
  description = "ARN of the SNS topic for alarms (if enabled)"
  value       = var.enable_cloudwatch_alarms && var.alarm_email != "" ? aws_sns_topic.lambda_alarms[0].arn : null
}

output "health_check_endpoint" {
  description = "Health check endpoint for the service"
  value       = "https://${var.dns_name}/members/health"
}