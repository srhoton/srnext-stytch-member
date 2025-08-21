# Stytch Member Management API Documentation

This directory contains the OpenAPI 3.0 specification for the Stytch Member Management Lambda API.

## Overview

The API provides full CRUD (Create, Read, Update, Delete) operations for managing members in Stytch B2B organizations. The Lambda function is designed to work behind an AWS Application Load Balancer (ALB) and handles ALB events.

## API Specification

The API specification is available in two formats:
- **YAML**: `openapi.yaml` - Human-readable format
- **JSON**: `openapi.json` - Machine-readable format for tools

## Endpoints

### Member Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/members` | Create a new member |
| GET | `/members/{member_id}` | Get a member by ID |
| PUT | `/members/{member_id}` | Update an existing member |
| DELETE | `/members/{member_id}` | Delete a member |
| POST | `/members/search` | Search for members |

## Authentication

The API supports two authentication methods:

1. **API Key**: Pass the API key in the `X-API-Key` header
2. **Bearer Token**: Pass a JWT token in the `Authorization: Bearer <token>` header

## Request/Response Format

- **Content-Type**: `application/json`
- **Accept**: `application/json`

## Common Response Codes

| Status Code | Description |
|------------|-------------|
| 200 | Success |
| 400 | Bad Request - Invalid input parameters |
| 401 | Unauthorized - Authentication failed |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource does not exist |
| 409 | Conflict - Resource already exists |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error |

## Error Response Format

All error responses follow this structure:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": "Additional error details"
  }
}
```

## Organization ID

The organization ID can be provided in two ways:
1. As a query parameter: `?organization_id=org-123`
2. As a header: `X-Organization-ID: org-123`

Query parameters take precedence if both are provided.

## Pagination

For search operations, pagination is supported using:
- `limit`: Maximum number of results (1-100, default: 20)
- `cursor`: Cursor from previous response for next page

## Examples

### Create a Member

```bash
curl -X POST https://api.example.com/members \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "organization_id": "org-123",
    "email_address": "user@example.com",
    "name": "John Doe",
    "roles": ["member"]
  }'
```

### Get a Member

```bash
curl -X GET "https://api.example.com/members/member-123?organization_id=org-123" \
  -H "X-API-Key: your-api-key"
```

### Update a Member

```bash
curl -X PUT https://api.example.com/members/member-123 \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "organization_id": "org-123",
    "member_id": "member-123",
    "name": "Jane Doe",
    "roles": ["admin", "member"]
  }'
```

### Delete a Member

```bash
curl -X DELETE "https://api.example.com/members/member-123?organization_id=org-123" \
  -H "X-API-Key: your-api-key"
```

### Search Members

```bash
curl -X POST https://api.example.com/members/search \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{
    "organization_ids": ["org-123"],
    "limit": 20,
    "query": {
      "role": "admin"
    }
  }'
```

## Metadata

The API supports two types of metadata:

1. **Trusted Metadata**: Can only be modified by the organization (server-side)
2. **Untrusted Metadata**: Can be modified by the member (client-side)

Both accept arbitrary JSON objects.

## MFA Support

Members can have Multi-Factor Authentication (MFA) configured:
- `mfa_phone_number`: Phone number for SMS-based MFA
- `mfa_enrolled`: Whether MFA is currently enrolled
- `default_mfa_method`: Default MFA method (`sms` or `totp`)

## Breakglass Accounts

The API supports creating and managing breakglass accounts using the `is_breakglass` flag. These are emergency access accounts with special privileges.

## Tools

### Viewing the Specification

You can use various tools to view and interact with the OpenAPI specification:

1. **Swagger UI**: Upload the spec to [Swagger Editor](https://editor.swagger.io/)
2. **Postman**: Import the spec into Postman for testing
3. **Insomnia**: Import the spec into Insomnia for testing
4. **ReDoc**: Generate beautiful documentation with [ReDoc](https://github.com/Redocly/redoc)

### Validation

To validate the OpenAPI specification:

```bash
npx @apidevtools/swagger-cli validate openapi.yaml
```

### Code Generation

You can generate client SDKs using OpenAPI Generator:

```bash
# Generate a Go client
openapi-generator generate -i openapi.yaml -g go -o ../client/go

# Generate a Python client
openapi-generator generate -i openapi.yaml -g python -o ../client/python

# Generate a TypeScript client
openapi-generator generate -i openapi.yaml -g typescript-axios -o ../client/typescript
```

## Integration with AWS API Gateway

This OpenAPI specification can be imported into AWS API Gateway to create a REST API that routes to the Lambda function. The specification includes all necessary request/response models and error responses.

## License

This API documentation follows the same license as the parent project.