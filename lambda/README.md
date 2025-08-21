# Stytch Member Management Lambda

AWS Lambda function for managing Stytch B2B organization members with full CRUD operations via AWS Application Load Balancer (ALB) events.

## Features

- **Create Member**: Add new members to Stytch organizations
- **Get Member**: Retrieve member details by ID or email
- **Update Member**: Modify member attributes and settings
- **Delete Member**: Remove members from organizations
- **Search Members**: Query members across organizations

## Architecture

- **Runtime**: AWS Lambda on `provided.al2023`
- **Language**: Go 1.24
- **SDK**: AWS SDK v2, Stytch Go SDK v12
- **Event Source**: AWS Application Load Balancer

## Project Structure

```
lambda/
├── cmd/lambda/         # Lambda entry point
├── internal/
│   ├── config/        # Configuration management
│   ├── handler/       # ALB event handler
│   ├── models/        # Request/response models
│   └── service/       # Stytch service implementation
├── pkg/errors/        # Error handling utilities
├── test/             # Test fixtures
├── Makefile          # Build automation
└── go.mod           # Go dependencies
```

## API Endpoints

### Create Member
```
POST /members
Content-Type: application/json

{
  "organization_id": "org-123",
  "email_address": "user@example.com",
  "name": "John Doe",
  "roles": ["member"],
  "trusted_metadata": {},
  "untrusted_metadata": {}
}
```

### Get Member
```
GET /members/{member_id}?organization_id=org-123
```

### Update Member
```
PUT /members/{member_id}
Content-Type: application/json

{
  "organization_id": "org-123",
  "name": "Jane Doe",
  "roles": ["admin"]
}
```

### Delete Member
```
DELETE /members/{member_id}?organization_id=org-123
```

### Search Members
```
POST /members/search
Content-Type: application/json

{
  "organization_ids": ["org-123"],
  "limit": 10,
  "cursor": "",
  "query": {}
}
```

## Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `STYTCH_PROJECT_ID` | Stytch project ID | Yes | - |
| `STYTCH_SECRET` | Stytch API secret | Yes | - |
| `STYTCH_ENV` | Stytch environment (`test` or `live`) | No | `test` |
| `LOG_LEVEL` | Logging level | No | `info` |
| `AWS_REGION` | AWS region | No | `us-east-1` |
| `REQUEST_TIMEOUT` | Request timeout in seconds | No | `30` |
| `MAX_RETRIES` | Maximum retry attempts | No | `3` |

## Building

### Prerequisites

- Go 1.24 or later
- Make
- AWS CLI (for deployment)

### Build Commands

```bash
# Install dependencies
make mod-download

# Format code
make fmt

# Run linters
make lint

# Run tests
make test

# Run tests with coverage
make test-coverage

# Build Lambda binary
make build

# Create deployment package
make package

# Run all checks and build
make all
```

## Testing

The project includes comprehensive unit tests with 80%+ code coverage requirement.

```bash
# Run tests
make test

# Run tests with coverage report
make test-coverage

# Verify coverage threshold (80%)
make verify-coverage
```

## Deployment

The Lambda function is packaged for deployment on AWS Lambda with the `provided.al2023` runtime.

```bash
# Create deployment package
make package

# Deploy using AWS CLI (example)
aws lambda update-function-code \
  --function-name stytch-member-lambda \
  --zip-file fileb://lambda-deployment.zip
```

## Security Considerations

- **Authentication**: Requires valid Stytch API credentials
- **CORS**: Configured for cross-origin requests
- **Input Validation**: All inputs are validated before processing
- **Error Handling**: Sensitive information is not exposed in error messages
- **TLS**: All Stytch API communications use HTTPS

## Error Handling

The Lambda returns standard HTTP status codes:

- `200 OK`: Successful operation
- `400 Bad Request`: Invalid input parameters
- `401 Unauthorized`: Authentication failure
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists
- `429 Too Many Requests`: Rate limited
- `500 Internal Server Error`: Unexpected error

## Development

### Local Testing

Set required environment variables:
```bash
export STYTCH_PROJECT_ID="your-project-id"
export STYTCH_SECRET="your-secret"
export STYTCH_ENV="test"
```

Run locally:
```bash
make run-local
```

### Code Quality

The project enforces code quality through:
- `golangci-lint` for comprehensive linting
- `gofmt`/`goimports` for formatting
- Unit tests with coverage requirements
- Type safety with no `interface{}` usage where possible

## License

[Your License Here]

## Support

For issues or questions, please contact the development team.