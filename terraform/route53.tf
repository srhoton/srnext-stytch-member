# Create DNS A record pointing to the ALB
resource "aws_route53_record" "main" {
  zone_id = data.aws_route53_zone.existing.zone_id
  name    = var.dns_name
  type    = "A"

  alias {
    name                   = data.aws_lb.existing.dns_name
    zone_id                = data.aws_lb.existing.zone_id
    evaluate_target_health = true
  }
}

# Optional: Create a health check for the endpoint
resource "aws_route53_health_check" "main" {
  count = var.enable_cloudwatch_alarms ? 1 : 0

  fqdn              = var.dns_name
  port              = 443
  type              = "HTTPS"
  resource_path     = "/members/health"
  failure_threshold = 3
  request_interval  = 30

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-health-check"
    }
  )
}

# CloudWatch alarm for Route53 health check
resource "aws_cloudwatch_metric_alarm" "route53_health" {
  count = var.enable_cloudwatch_alarms ? 1 : 0

  alarm_name          = "${local.name_prefix}-route53-health"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = 2
  metric_name         = "HealthCheckStatus"
  namespace           = "AWS/Route53"
  period              = 60
  statistic           = "Minimum"
  threshold           = 1
  alarm_description   = "This metric monitors Route53 health check status"
  treat_missing_data  = "breaching"

  dimensions = {
    HealthCheckId = aws_route53_health_check.main[0].id
  }

  alarm_actions = var.alarm_email != "" ? [aws_sns_topic.lambda_alarms[0].arn] : []

  tags = local.common_tags
}