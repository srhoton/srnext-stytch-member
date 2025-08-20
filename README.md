# Stytch Member Lambda Service

AWS Lambda service for managing Stytch members with full CRUD operations through API Gateway V2 events.

## Project Structure

```
srnextstytchmember/
├── src/
│   ├── main/java/com/steverhoton/poc/
│   │   ├── lambda/           # Lambda handler
│   │   ├── service/          # Business logic
│   │   ├── model/            # Request/Response models
│   │   ├── exception/        # Custom exceptions
│   │   └── stytch/mock/      # Mock Stytch models
│   └── test/                 # Unit tests
├── build.gradle              # Build configuration
└── settings.gradle           # Project settings
```

## Features

- **Complete CRUD Operations**: Create, Read, Update, Delete, and Reactivate members
- **API Gateway V2 Support**: Handles HTTP API events
- **Error Handling**: Comprehensive error handling with specific exceptions
- **Security**: Input validation and sanitization
- **Testing**: 86% test coverage with unit tests
- **Code Quality**: Spotless formatting and linting

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /members | Create a new member |
| GET | /members | Search members |
| GET | /members/{id} | Get member by ID |
| PUT | /members/{id} | Update member |
| DELETE | /members/{id} | Delete member |
| POST | /members/{id}/reactivate | Reactivate member |

## Environment Variables

The Lambda requires the following environment variables:

- `STYTCH_PROJECT_ID`: Your Stytch project ID
- `STYTCH_SECRET`: Your Stytch API secret
- `STYTCH_ORGANIZATION_ID`: Organization ID (optional, defaults to "default-org")

## Building

```bash
# Build the project
./gradlew build

# Run tests
./gradlew test

# Check test coverage
./gradlew jacocoTestReport

# Apply formatting
./gradlew spotlessApply

# Create Lambda deployment package
./gradlew buildZip
```

## Deployment Package

The Lambda deployment package is created at `build/distributions/srnextstytchmember.zip`

## Lambda Configuration

- **Runtime**: Java 21
- **Handler**: `com.steverhoton.poc.lambda.MemberLambdaHandler::handleRequest`
- **Memory**: 512 MB (recommended)
- **Timeout**: 30 seconds (recommended)

## Request/Response Examples

### Create Member
```json
POST /members
{
  "email_address": "user@example.com",
  "name": "John Doe",
  "trusted_metadata": {},
  "untrusted_metadata": {},
  "create_member_as_pending": false,
  "is_breakglass": false,
  "mfa_phone_number": "+1234567890",
  "mfa_enrolled": false
}
```

### Search Members
```
GET /members?email=user@example.com&limit=10
```

### Update Member
```json
PUT /members/{memberId}
{
  "name": "Jane Doe",
  "mfa_enrolled": true
}
```

## Error Responses

All error responses follow this format:
```json
{
  "error": "Error Type",
  "message": "Detailed error message",
  "timestamp": "2025-01-20T12:00:00Z",
  "path": "/members/123"
}
```

## Testing

The project includes comprehensive unit tests:

- Service layer tests with mock Stytch implementation
- Lambda handler tests with various scenarios
- Model and exception tests
- **Test Coverage**: 86%

## Important Notes

1. **Mock Implementation**: This implementation uses mock Stytch models for demonstration. In production, integrate with the actual Stytch SDK v7.26.0 or later.

2. **Security**: Always use HTTPS for API Gateway endpoints and never expose sensitive credentials.

3. **Logging**: Logs are written using SLF4J and can be viewed in CloudWatch Logs.

## Dependencies

- AWS Lambda Java Core
- AWS Lambda Java Events
- AWS SDK v2
- Jackson for JSON processing
- Apache Commons Lang3
- SLF4J/Log4j2 for logging
- JUnit 5, Mockito, AssertJ for testing

## Development

1. Ensure Java 21 is installed
2. Clone the repository
3. Set up environment variables for local testing
4. Run tests: `./gradlew test`
5. Build: `./gradlew build`

## License

This is a proof of concept project.
