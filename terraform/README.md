# Terraform Infrastructure for Stytch Member Management Lambda

This Terraform configuration deploys the Stytch Member Management Lambda function and all associated AWS infrastructure.

## Architecture

The infrastructure includes:

- **Lambda Function**: Go-based function with `provided.al2023` runtime
- **DynamoDB Table**: For storing member data with email-based GSI
- **ALB Integration**: Target group and listener rules for existing ALB
- **Route53**: DNS configuration for the service endpoint
- **CloudWatch**: Logging, monitoring, alarms, and dashboard
- **IAM**: Roles and policies for Lambda execution
- **SNS**: Topic for alarm notifications (optional)

## Prerequisites

- Terraform >= 1.5.0
- AWS CLI configured with appropriate credentials
- Go 1.24+ installed locally (for Lambda build)
- Make installed (for build automation)
- Access to the S3 backend bucket: `steve-rhoton-tfstate`
- Existing ALB and Route53 hosted zone

## Configuration

1. Copy the example variables file:
```bash
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your configuration:
   - **Required**: Set `stytch_project_id` and `stytch_secret`
   - **Optional**: Adjust other parameters as needed

## Deployment

### Initialize Terraform

```bash
terraform init
```

This will configure the S3 backend and download required providers.

### Format and Validate

```bash
# Format Terraform files
terraform fmt -recursive

# Validate configuration
terraform validate

# Run tflint for additional checks
tflint --init
tflint
```

### Plan Changes

```bash
terraform plan -out=tfplan
```

Review the planned changes carefully.

### Apply Changes

```bash
terraform apply tfplan
```

Or directly:
```bash
terraform apply
```

### Destroy Infrastructure

```bash
terraform destroy
```

## State Management

Terraform state is stored in S3:
- **Bucket**: `steve-rhoton-tfstate`
- **Key**: `srnext-stytch-member/terraform.tfstate`
- **Region**: `us-west-2`

## Resources Created

### Core Infrastructure

| Resource | Name Pattern | Description |
|----------|-------------|-------------|
| Lambda Function | `{project}-{env}-lambda` | Main function for API operations |
| DynamoDB Table | `{project}-{env}-members` | Member data storage |
| CloudWatch Log Group | `/aws/lambda/{function-name}` | Function logs |
| IAM Role | `{function-name}-execution-role` | Lambda execution role |

### ALB Configuration

| Resource | Description |
|----------|-------------|
| Target Group | Lambda target group for ALB |
| Listener Rule | Routes `/members/*` paths to Lambda |
| Lambda Permission | Allows ALB to invoke Lambda |

### DNS Configuration

| Resource | Description |
|----------|-------------|
| Route53 A Record | Points DNS name to ALB |
| Health Check | Monitors endpoint availability (optional) |

### Monitoring

| Alarm | Threshold | Description |
|-------|-----------|-------------|
| Error Rate | 5 errors/min | Lambda function errors |
| Throttles | 10/min | Lambda throttling |
| Duration | 80% of timeout | Execution duration |
| Concurrent Executions | 90% of reserved | Concurrency usage |

## Environment Variables

The Lambda function receives these environment variables:

| Variable | Description |
|----------|-------------|
| `STYTCH_PROJECT_ID` | Stytch project identifier |
| `STYTCH_SECRET` | Stytch API secret |
| `STYTCH_ENV` | Stytch environment (test/live) |
| `AWS_REGION` | AWS region |
| `DYNAMODB_TABLE` | DynamoDB table name |
| `LOG_LEVEL` | Logging level |
| `REQUEST_TIMEOUT` | Request timeout in seconds |
| `MAX_RETRIES` | Maximum retry attempts |

## Customization

### VPC Configuration

To deploy the Lambda in a VPC:

```hcl
vpc_subnet_ids         = ["subnet-xxx", "subnet-yyy"]
vpc_security_group_ids = ["sg-xxx"]
```

### Provisioned DynamoDB

To use provisioned capacity instead of on-demand:

```hcl
dynamodb_billing_mode  = "PROVISIONED"
dynamodb_read_capacity = 10
dynamodb_write_capacity = 10
```

### Custom Memory and Timeout

```hcl
lambda_memory_size = 1024  # MB
lambda_timeout     = 60    # seconds
```

## Monitoring

### CloudWatch Dashboard

Access the dashboard at:
```
https://us-west-2.console.aws.amazon.com/cloudwatch/home?region=us-west-2#dashboards:name={project}-{env}
```

### Alarms

If `enable_cloudwatch_alarms = true` and `alarm_email` is set:
- Email notifications will be sent for alarm breaches
- Confirm the SNS subscription email to receive notifications

### Logs

View Lambda logs in CloudWatch:
```bash
aws logs tail /aws/lambda/{function-name} --follow
```

## API Endpoints

After deployment, the API will be available at:

```
https://srnext-stytch-member.sb.int.fullbayapi.com
```

### Endpoints

- `POST /members` - Create member
- `GET /members/{id}` - Get member
- `PUT /members/{id}` - Update member
- `DELETE /members/{id}` - Delete member
- `POST /members/search` - Search members
- `GET /members/health` - Health check

## Troubleshooting

### Build Issues

If the Lambda build fails:

1. Ensure Go 1.24+ is installed:
```bash
go version
```

2. Manually build the Lambda:
```bash
cd ../lambda
make clean build
```

### Permission Issues

If you get permission errors:

1. Verify AWS credentials:
```bash
aws sts get-caller-identity
```

2. Check IAM permissions for:
   - Lambda operations
   - DynamoDB operations
   - ALB modifications
   - Route53 changes
   - CloudWatch operations

### State Lock Issues

If Terraform state is locked:

1. Wait for other operations to complete
2. If stuck, manually check S3 for lock files
3. Use `terraform force-unlock` with caution

## Cost Considerations

### Pay-Per-Use Resources

- Lambda: Charged per invocation and GB-seconds
- DynamoDB (On-Demand): Charged per read/write request
- CloudWatch Logs: Charged per GB ingested and stored
- ALB: Already exists, no additional charge for rules

### Fixed Costs

- Route53 Health Check: ~$0.50/month per health check
- CloudWatch Alarms: ~$0.10/month per alarm
- SNS: ~$0.50 per million notifications

### Cost Optimization

- Adjust `log_retention_days` to reduce storage costs
- Use `lambda_reserved_concurrent_executions` to control costs
- Monitor DynamoDB usage and switch to provisioned if predictable

## Security Best Practices

1. **Secrets Management**:
   - Never commit `terraform.tfvars` with secrets
   - Consider using AWS Secrets Manager for Stytch credentials
   - Use environment-specific secrets

2. **IAM Permissions**:
   - Lambda role follows least privilege principle
   - Separate roles for different environments

3. **Network Security**:
   - Consider VPC deployment for additional isolation
   - Use security groups to restrict access

4. **Encryption**:
   - DynamoDB encryption at rest is enabled
   - CloudWatch Logs are encrypted
   - Use HTTPS for all API communications

## Backup and Recovery

### DynamoDB

- Point-in-time recovery is enabled
- Manual backups can be created:
```bash
aws dynamodb create-backup \
  --table-name {table-name} \
  --backup-name {backup-name}
```

### Lambda Versions

- Each deployment creates a new version
- Use aliases for rollback capability
- Previous versions are retained

## Maintenance

### Updating the Lambda

1. Make code changes in `../lambda/`
2. Run `terraform apply` to deploy

### Updating Dependencies

```bash
# Update Terraform providers
terraform init -upgrade

# Update Go dependencies
cd ../lambda
go get -u ./...
go mod tidy
```

### Monitoring Costs

```bash
# View cost breakdown
aws ce get-cost-and-usage \
  --time-period Start=2024-01-01,End=2024-01-31 \
  --granularity MONTHLY \
  --metrics "BlendedCost" \
  --group-by Type=DIMENSION,Key=SERVICE
```

## Support

For issues or questions:
1. Check CloudWatch Logs for errors
2. Review CloudWatch metrics and alarms
3. Verify IAM permissions
4. Check ALB target health
5. Test DNS resolution

## License

[Your License Here]