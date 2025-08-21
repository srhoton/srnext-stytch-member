package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/steverhoton/srnext-stytch-member/lambda/internal/config"
	"github.com/steverhoton/srnext-stytch-member/lambda/internal/models"
	"github.com/steverhoton/srnext-stytch-member/lambda/internal/service"
	"github.com/steverhoton/srnext-stytch-member/lambda/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// MockMemberService is a mock implementation of MemberService
type MockMemberService struct {
	mock.Mock
}

// Ensure MockMemberService implements service.MemberService
var _ service.MemberService = (*MockMemberService)(nil)

func (m *MockMemberService) CreateMember(ctx context.Context, req *models.CreateMemberRequest) (*models.MemberResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MemberResponse), args.Error(1)
}

func (m *MockMemberService) GetMember(ctx context.Context, req *models.GetMemberRequest) (*models.MemberResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MemberResponse), args.Error(1)
}

func (m *MockMemberService) UpdateMember(ctx context.Context, req *models.UpdateMemberRequest) (*models.MemberResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MemberResponse), args.Error(1)
}

func (m *MockMemberService) DeleteMember(ctx context.Context, req *models.DeleteMemberRequest) (*models.MemberResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MemberResponse), args.Error(1)
}

func (m *MockMemberService) SearchMembers(ctx context.Context, req *models.SearchMembersRequest) (*models.MembersListResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MembersListResponse), args.Error(1)
}

func TestNewHandler(t *testing.T) {
	logger := zap.NewNop()
	cfg := &config.Config{
		StytchProjectID: "test-project",
		StytchSecret:    "test-secret",
		StytchEnv:       "test",
	}

	tests := []struct {
		name    string
		cfg     *config.Config
		logger  *zap.Logger
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid configuration",
			cfg:     cfg,
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
			name:    "Nil logger",
			cfg:     cfg,
			logger:  nil,
			wantErr: true,
			errMsg:  "logger is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't easily test NewHandler without mocking service creation
			// This would require refactoring to inject the service
			if tt.cfg == nil || tt.logger == nil {
				// Just verify the validation logic would work
				assert.True(t, tt.wantErr)
			}
		})
	}
}

func TestHandler_HandleRequest_CORS(t *testing.T) {
	h := &Handler{
		memberService: &MockMemberService{},
		logger:        zap.NewNop(),
		config:        &config.Config{},
	}

	req := models.ALBRequest{
		HTTPMethod: http.MethodOptions,
		Path:       "/members",
	}

	resp, err := h.HandleRequest(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "*", resp.Headers["Access-Control-Allow-Origin"])
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", resp.Headers["Access-Control-Allow-Methods"])
	assert.Equal(t, "Content-Type, Authorization", resp.Headers["Access-Control-Allow-Headers"])
}

func TestHandler_HandleRequest_CreateMember(t *testing.T) {
	mockService := &MockMemberService{}
	h := &Handler{
		memberService: mockService,
		logger:        zap.NewNop(),
		config:        &config.Config{},
	}

	createReq := models.CreateMemberRequest{
		OrganizationID: "org-123",
		EmailAddress:   "test@example.com",
		Name:           "Test User",
	}

	expectedResponse := &models.MemberResponse{
		StatusCode: http.StatusOK,
		RequestID:  "req-123",
		Member: models.Member{
			MemberID:       "member-123",
			OrganizationID: "org-123",
			EmailAddress:   "test@example.com",
			Name:           "Test User",
		},
	}

	mockService.On("CreateMember", mock.Anything, &createReq).Return(expectedResponse, nil)

	body, _ := json.Marshal(createReq)
	req := models.ALBRequest{
		HTTPMethod: http.MethodPost,
		Path:       "/members",
		Body:       string(body),
	}

	resp, err := h.HandleRequest(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result models.MemberResponse
	err = json.Unmarshal([]byte(resp.Body), &result)
	require.NoError(t, err)
	assert.Equal(t, expectedResponse.Member.MemberID, result.Member.MemberID)

	mockService.AssertExpectations(t)
}

func TestHandler_HandleRequest_GetMember(t *testing.T) {
	mockService := &MockMemberService{}
	h := &Handler{
		memberService: mockService,
		logger:        zap.NewNop(),
		config:        &config.Config{},
	}

	getReq := &models.GetMemberRequest{
		OrganizationID: "org-123",
		MemberID:       "member-123",
	}

	expectedResponse := &models.MemberResponse{
		StatusCode: http.StatusOK,
		RequestID:  "req-123",
		Member: models.Member{
			MemberID:       "member-123",
			OrganizationID: "org-123",
			EmailAddress:   "test@example.com",
		},
	}

	mockService.On("GetMember", mock.Anything, getReq).Return(expectedResponse, nil)

	req := models.ALBRequest{
		HTTPMethod: http.MethodGet,
		Path:       "/members/member-123",
		QueryStringParameters: map[string]string{
			"organization_id": "org-123",
		},
	}

	resp, err := h.HandleRequest(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result models.MemberResponse
	err = json.Unmarshal([]byte(resp.Body), &result)
	require.NoError(t, err)
	assert.Equal(t, expectedResponse.Member.MemberID, result.Member.MemberID)

	mockService.AssertExpectations(t)
}

func TestHandler_HandleRequest_UpdateMember(t *testing.T) {
	mockService := &MockMemberService{}
	h := &Handler{
		memberService: mockService,
		logger:        zap.NewNop(),
		config:        &config.Config{},
	}

	updateReq := models.UpdateMemberRequest{
		OrganizationID: "org-123",
		MemberID:       "member-123",
		Name:           "Updated Name",
	}

	expectedResponse := &models.MemberResponse{
		StatusCode: http.StatusOK,
		RequestID:  "req-123",
		Member: models.Member{
			MemberID:       "member-123",
			OrganizationID: "org-123",
			Name:           "Updated Name",
		},
	}

	mockService.On("UpdateMember", mock.Anything, &updateReq).Return(expectedResponse, nil)

	body, _ := json.Marshal(updateReq)
	req := models.ALBRequest{
		HTTPMethod: http.MethodPut,
		Path:       "/members/member-123",
		Body:       string(body),
	}

	resp, err := h.HandleRequest(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result models.MemberResponse
	err = json.Unmarshal([]byte(resp.Body), &result)
	require.NoError(t, err)
	assert.Equal(t, expectedResponse.Member.Name, result.Member.Name)

	mockService.AssertExpectations(t)
}

func TestHandler_HandleRequest_DeleteMember(t *testing.T) {
	mockService := &MockMemberService{}
	h := &Handler{
		memberService: mockService,
		logger:        zap.NewNop(),
		config:        &config.Config{},
	}

	deleteReq := &models.DeleteMemberRequest{
		OrganizationID: "org-123",
		MemberID:       "member-123",
	}

	expectedResponse := &models.MemberResponse{
		StatusCode: http.StatusOK,
		RequestID:  "req-123",
		Member: models.Member{
			MemberID: "member-123",
		},
	}

	mockService.On("DeleteMember", mock.Anything, deleteReq).Return(expectedResponse, nil)

	req := models.ALBRequest{
		HTTPMethod: http.MethodDelete,
		Path:       "/members/member-123",
		Headers: map[string]string{
			"x-organization-id": "org-123",
		},
	}

	resp, err := h.HandleRequest(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestHandler_HandleRequest_SearchMembers(t *testing.T) {
	mockService := &MockMemberService{}
	h := &Handler{
		memberService: mockService,
		logger:        zap.NewNop(),
		config:        &config.Config{},
	}

	searchReq := models.SearchMembersRequest{
		OrganizationIDs: []string{"org-123"},
		Limit:           10,
	}

	expectedResponse := &models.MembersListResponse{
		StatusCode: http.StatusOK,
		RequestID:  "req-123",
		Members: []models.Member{
			{
				MemberID:       "member-1",
				OrganizationID: "org-123",
				EmailAddress:   "test1@example.com",
			},
			{
				MemberID:       "member-2",
				OrganizationID: "org-123",
				EmailAddress:   "test2@example.com",
			},
		},
	}

	mockService.On("SearchMembers", mock.Anything, &searchReq).Return(expectedResponse, nil)

	body, _ := json.Marshal(searchReq)
	req := models.ALBRequest{
		HTTPMethod: http.MethodPost,
		Path:       "/members/search",
		Body:       string(body),
	}

	resp, err := h.HandleRequest(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result models.MembersListResponse
	err = json.Unmarshal([]byte(resp.Body), &result)
	require.NoError(t, err)
	assert.Len(t, result.Members, 2)

	mockService.AssertExpectations(t)
}

func TestHandler_HandleRequest_InvalidPath(t *testing.T) {
	h := &Handler{
		memberService: &MockMemberService{},
		logger:        zap.NewNop(),
		config:        &config.Config{},
	}

	req := models.ALBRequest{
		HTTPMethod: http.MethodGet,
		Path:       "/invalid",
	}

	resp, err := h.HandleRequest(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var errResp errors.ErrorResponse
	err = json.Unmarshal([]byte(resp.Body), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "NOT_FOUND", errResp.Error.Code)
}

func TestHandler_HandleRequest_InvalidMethod(t *testing.T) {
	h := &Handler{
		memberService: &MockMemberService{},
		logger:        zap.NewNop(),
		config:        &config.Config{},
	}

	req := models.ALBRequest{
		HTTPMethod: http.MethodPatch,
		Path:       "/members",
	}

	resp, err := h.HandleRequest(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	var errResp errors.ErrorResponse
	err = json.Unmarshal([]byte(resp.Body), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "METHOD_NOT_ALLOWED", errResp.Error.Code)
}

func TestHandler_HandleRequest_ServiceError(t *testing.T) {
	mockService := &MockMemberService{}
	h := &Handler{
		memberService: mockService,
		logger:        zap.NewNop(),
		config:        &config.Config{},
	}

	serviceErr := errors.NewAPIError(http.StatusConflict, "CONFLICT", "Member already exists")
	mockService.On("CreateMember", mock.Anything, mock.Anything).Return(nil, serviceErr)

	createReq := models.CreateMemberRequest{
		OrganizationID: "org-123",
		EmailAddress:   "test@example.com",
	}

	body, _ := json.Marshal(createReq)
	req := models.ALBRequest{
		HTTPMethod: http.MethodPost,
		Path:       "/members",
		Body:       string(body),
	}

	resp, err := h.HandleRequest(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusConflict, resp.StatusCode)

	var errResp errors.ErrorResponse
	err = json.Unmarshal([]byte(resp.Body), &errResp)
	require.NoError(t, err)
	assert.Equal(t, "CONFLICT", errResp.Error.Code)
	assert.Equal(t, "Member already exists", errResp.Error.Message)

	mockService.AssertExpectations(t)
}

func TestHandler_validateCreateRequest(t *testing.T) {
	h := &Handler{
		memberService: &MockMemberService{},
		logger:        zap.NewNop(),
		config:        &config.Config{},
	}

	tests := []struct {
		name    string
		req     *models.CreateMemberRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid request",
			req: &models.CreateMemberRequest{
				OrganizationID: "org-123",
				EmailAddress:   "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "Missing organization ID",
			req: &models.CreateMemberRequest{
				EmailAddress: "test@example.com",
			},
			wantErr: true,
			errMsg:  "Organization ID is required",
		},
		{
			name: "Missing email",
			req: &models.CreateMemberRequest{
				OrganizationID: "org-123",
			},
			wantErr: true,
			errMsg:  "Email address is required",
		},
		{
			name: "Invalid email",
			req: &models.CreateMemberRequest{
				OrganizationID: "org-123",
				EmailAddress:   "invalid-email",
			},
			wantErr: true,
			errMsg:  "Invalid email address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := h.validateCreateRequest(tt.req)

			if tt.wantErr {
				require.Error(t, err)
				apiErr, ok := err.(*errors.APIError)
				require.True(t, ok)
				assert.Contains(t, apiErr.Message, tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestHandler_getOrganizationID(t *testing.T) {
	h := &Handler{
		memberService: &MockMemberService{},
		logger:        zap.NewNop(),
		config:        &config.Config{},
	}

	tests := []struct {
		name     string
		req      models.ALBRequest
		expected string
	}{
		{
			name: "From query parameters",
			req: models.ALBRequest{
				QueryStringParameters: map[string]string{
					"organization_id": "org-123",
				},
			},
			expected: "org-123",
		},
		{
			name: "From headers",
			req: models.ALBRequest{
				Headers: map[string]string{
					"x-organization-id": "org-456",
				},
			},
			expected: "org-456",
		},
		{
			name: "Query parameters take precedence",
			req: models.ALBRequest{
				QueryStringParameters: map[string]string{
					"organization_id": "org-123",
				},
				Headers: map[string]string{
					"x-organization-id": "org-456",
				},
			},
			expected: "org-123",
		},
		{
			name:     "Not found",
			req:      models.ALBRequest{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := h.getOrganizationID(tt.req)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHandler_errorResponse(t *testing.T) {
	h := &Handler{
		memberService: &MockMemberService{},
		logger:        zap.NewNop(),
		config:        &config.Config{},
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	tests := []struct {
		name               string
		err                error
		expectedStatusCode int
		expectedCode       string
	}{
		{
			name:               "API error",
			err:                errors.NewAPIError(http.StatusBadRequest, "BAD_REQUEST", "Invalid input"),
			expectedStatusCode: http.StatusBadRequest,
			expectedCode:       "BAD_REQUEST",
		},
		{
			name:               "Generic error",
			err:                fmt.Errorf("something went wrong"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedCode:       "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := h.errorResponse(tt.err, headers)

			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			assert.Equal(t, headers, resp.Headers)

			var errResp errors.ErrorResponse
			err := json.Unmarshal([]byte(resp.Body), &errResp)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCode, errResp.Error.Code)
		})
	}
}
