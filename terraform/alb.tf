# Target group for Lambda
resource "aws_lb_target_group" "lambda" {
  name        = local.target_group_name
  target_type = "lambda"

  # Health check configuration for Lambda target
  health_check {
    enabled             = true
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 10
    interval            = 30
    path                = "/members/health"
    matcher             = "200"
  }

  tags = merge(
    local.common_tags,
    {
      Name = local.target_group_name
    }
  )
}

# Attach Lambda to target group
resource "aws_lb_target_group_attachment" "lambda" {
  target_group_arn = aws_lb_target_group.lambda.arn
  target_id        = aws_lambda_function.main.arn

  depends_on = [aws_lambda_permission.alb_invoke]
}

# ALB listener rule for routing to Lambda
resource "aws_lb_listener_rule" "lambda" {
  listener_arn = data.aws_lb_listener.existing.arn
  priority     = var.alb_rule_priority

  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.lambda.arn
  }

  # Path-based routing
  condition {
    path_pattern {
      values = var.alb_path_patterns
    }
  }

  # Host-based routing
  condition {
    host_header {
      values = [var.dns_name]
    }
  }

  tags = merge(
    local.common_tags,
    {
      Name = "${local.name_prefix}-listener-rule"
    }
  )
}