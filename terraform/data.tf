data "aws_caller_identity" "current" {}

# Stytch credentials from Secrets Manager
data "aws_secretsmanager_secret" "stytch_credentials" {
  name = "srnext/stytchCredentials"
}

data "aws_secretsmanager_secret_version" "stytch_credentials" {
  secret_id = data.aws_secretsmanager_secret.stytch_credentials.id
}

# Existing ALB
data "aws_lb" "existing" {
  arn = var.alb_arn
}

# Existing ALB listener
data "aws_lb_listener" "existing" {
  load_balancer_arn = data.aws_lb.existing.arn
  port              = var.alb_listener_port
}

# Existing Route53 hosted zone
data "aws_route53_zone" "existing" {
  name         = var.dns_zone_name
  private_zone = false
}