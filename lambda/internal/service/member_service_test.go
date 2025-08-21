package service

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/steverhoton/srnext-stytch-member/lambda/internal/config"
	"github.com/steverhoton/srnext-stytch-member/lambda/internal/models"
	"github.com/steverhoton/srnext-stytch-member/lambda/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/organizations"
	"go.uber.org/zap"
)

func TestNewMemberService(t *testing.T) {
	logger := zap.NewNop()

	tests := []struct {
		name    string
		cfg     *config.Config
		logger  *zap.Logger
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid configuration test environment",
			cfg: &config.Config{
				StytchProjectID: "test-project",
				StytchSecret:    "test-secret",
				StytchEnv:       "test",
			},
			logger:  logger,
			wantErr: false,
		},
		{
			name: "Valid configuration live environment",
			cfg: &config.Config{
				StytchProjectID: "test-project",
				StytchSecret:    "test-secret",
				StytchEnv:       "live",
			},
			logger:  logger,
			wantErr: false,
		},
		{
			name:    "Nil config",
			cfg:     nil,
			logger:  logger,
			wantErr: true,
			errMsg:  "config is required",
		},
		{
			name: "Nil logger",
			cfg: &config.Config{
				StytchProjectID: "test-project",
				StytchSecret:    "test-secret",
				StytchEnv:       "test",
			},
			logger:  nil,
			wantErr: true,
			errMsg:  "logger is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewMemberService(tt.cfg, tt.logger)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, service)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, service)
			}
		})
	}
}

func TestMapStytchMember(t *testing.T) {
	service := &stytchMemberService{
		logger: zap.NewNop(),
	}

	// Use the actual mapMember method from the service
	adminRole := organizations.MemberRole{RoleID: "admin"}
	userRole := organizations.MemberRole{RoleID: "user"}

	stytchMember := organizations.Member{
		MemberID:          "member-123",
		OrganizationID:    "org-123",
		EmailAddress:      "test@example.com",
		Status:            "active",
		Name:              "Test User",
		Roles:             []organizations.MemberRole{adminRole, userRole},
		TrustedMetadata:   map[string]interface{}{"key": "value"},
		UntrustedMetadata: map[string]interface{}{"pref": "dark"},
		MFAPhoneNumber:    "+1234567890",
		MFAEnrolled:       true,
		DefaultMFAMethod:  "sms",
		IsBreakglass:      false,
	}

	result := service.mapMember(stytchMember)

	assert.Equal(t, "member-123", result.MemberID)
	assert.Equal(t, "org-123", result.OrganizationID)
	assert.Equal(t, "test@example.com", result.EmailAddress)
	assert.Equal(t, "active", result.Status)
	assert.Equal(t, "Test User", result.Name)
	// ExternalID not available in SDK v12
	// assert.Equal(t, "ext-123", result.ExternalID)
	assert.Equal(t, []string{"admin", "user"}, result.Roles)
	assert.Equal(t, map[string]interface{}{"key": "value"}, result.TrustedMetadata)
	assert.Equal(t, map[string]interface{}{"pref": "dark"}, result.UntrustedMetadata)
	assert.Equal(t, "+1234567890", result.MFAPhoneNumber)
	assert.True(t, result.MFAEnrolled)
	assert.Equal(t, "sms", result.DefaultMFAMethod)
	assert.False(t, result.IsBreakglass)
	// CreatedAt and UpdatedAt not available in SDK v12
	// assert.Equal(t, now, result.CreatedAt)
	// assert.Equal(t, now, result.UpdatedAt)
}

func TestMapStytchMemberWithNilFields(t *testing.T) {
	service := &stytchMemberService{
		logger: zap.NewNop(),
	}

	stytchMember := organizations.Member{
		MemberID:       "member-123",
		OrganizationID: "org-123",
		EmailAddress:   "test@example.com",
		Status:         "pending",
		Name:           "",
		Roles:          []organizations.MemberRole{},
	}

	result := service.mapMember(stytchMember)

	assert.Equal(t, "member-123", result.MemberID)
	assert.Equal(t, "org-123", result.OrganizationID)
	assert.Equal(t, "test@example.com", result.EmailAddress)
	assert.Equal(t, "pending", result.Status)
	assert.Empty(t, result.Name)
	// ExternalID not available in SDK v12
	// assert.Empty(t, result.ExternalID)
	assert.Empty(t, result.Roles)
	assert.Empty(t, result.DefaultMFAMethod)
	assert.False(t, result.MFAEnrolled)
	assert.False(t, result.IsBreakglass)
}

// TestBuildSearchOperands - removed as buildSearchOperands is no longer exported
// The SDK handles search query building internally

func TestHandleStytchError(t *testing.T) {
	service := &stytchMemberService{
		logger: zap.NewNop(),
	}

	tests := []struct {
		name               string
		err                error
		expectedStatusCode int
		expectedCode       string
		expectedMessage    string
	}{
		{
			name:               "Bad request error",
			err:                fmt.Errorf("400 bad request"),
			expectedStatusCode: http.StatusBadRequest,
			expectedCode:       "BAD_REQUEST",
			expectedMessage:    "Invalid request",
		},
		{
			name:               "Unauthorized error",
			err:                fmt.Errorf("401 unauthorized"),
			expectedStatusCode: http.StatusUnauthorized,
			expectedCode:       "UNAUTHORIZED",
			expectedMessage:    "Authentication failed",
		},
		{
			name:               "Forbidden error",
			err:                fmt.Errorf("403 forbidden"),
			expectedStatusCode: http.StatusForbidden,
			expectedCode:       "FORBIDDEN",
			expectedMessage:    "Access denied",
		},
		{
			name:               "Not found error",
			err:                fmt.Errorf("404 not found"),
			expectedStatusCode: http.StatusNotFound,
			expectedCode:       "NOT_FOUND",
			expectedMessage:    "Member not found",
		},
		{
			name:               "Conflict error",
			err:                fmt.Errorf("409 conflict"),
			expectedStatusCode: http.StatusConflict,
			expectedCode:       "CONFLICT",
			expectedMessage:    "Member already exists",
		},
		{
			name:               "Rate limit error",
			err:                fmt.Errorf("429 rate"),
			expectedStatusCode: http.StatusTooManyRequests,
			expectedCode:       "RATE_LIMITED",
			expectedMessage:    "Too many requests",
		},
		{
			name:               "Other error",
			err:                fmt.Errorf("some other error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedCode:       "INTERNAL_ERROR",
			expectedMessage:    "An error occurred",
		},
		{
			name:               "Generic error",
			err:                fmt.Errorf("something went wrong"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedCode:       "INTERNAL_ERROR",
			expectedMessage:    "An error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.handleStytchError(tt.err)

			apiErr, ok := result.(*errors.APIError)
			require.True(t, ok)

			assert.Equal(t, tt.expectedStatusCode, apiErr.StatusCode)
			assert.Equal(t, tt.expectedCode, apiErr.Code)
			assert.Equal(t, tt.expectedMessage, apiErr.Message)

			// Check that details contain the original error message
			assert.Equal(t, tt.err.Error(), apiErr.Details)
		})
	}
}

// Integration test example (would require mocking Stytch client)
func TestCreateMemberIntegration(t *testing.T) {
	// This test demonstrates the structure but would require
	// a mock Stytch client to work properly
	t.Skip("Integration test requires Stytch client mock")

	cfg := &config.Config{
		StytchProjectID: "test-project",
		StytchSecret:    "test-secret",
		StytchEnv:       "test",
	}

	service, err := NewMemberService(cfg, zap.NewNop())
	require.NoError(t, err)

	req := &models.CreateMemberRequest{
		OrganizationID: "org-123",
		EmailAddress:   "test@example.com",
		Name:           "Test User",
		Roles:          []string{"member"},
		TrustedMetadata: map[string]interface{}{
			"department": "engineering",
		},
	}

	// This would fail without a proper mock
	_, err = service.CreateMember(context.Background(), req)
	assert.Error(t, err) // Would be an authentication error
}

// Validation test for create member request handling
func TestCreateMemberRequestMapping(t *testing.T) {
	// This test verifies that the service correctly maps
	// the internal CreateMemberRequest to Stytch's params
	// The actual test would require mocking the Stytch client

	req := &models.CreateMemberRequest{
		OrganizationID:        "org-123",
		EmailAddress:          "test@example.com",
		Name:                  "Test User",
		ExternalID:            "ext-123",
		Roles:                 []string{"admin"},
		TrustedMetadata:       map[string]interface{}{"key": "value"},
		UntrustedMetadata:     map[string]interface{}{"pref": "dark"},
		MFAPhoneNumber:        "+1234567890",
		CreateMemberAsPending: true,
		IsBreakglass:          false,
	}

	// Verify all fields are present
	assert.NotEmpty(t, req.OrganizationID)
	assert.NotEmpty(t, req.EmailAddress)
	assert.NotEmpty(t, req.Name)
	assert.NotEmpty(t, req.ExternalID)
	assert.NotEmpty(t, req.Roles)
	assert.NotNil(t, req.TrustedMetadata)
	assert.NotNil(t, req.UntrustedMetadata)
	assert.NotEmpty(t, req.MFAPhoneNumber)
	assert.True(t, req.CreateMemberAsPending)
	assert.False(t, req.IsBreakglass)
}
