package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/steverhoton/srnext-stytch-member/lambda/internal/config"
	"github.com/steverhoton/srnext-stytch-member/lambda/internal/models"
	"github.com/steverhoton/srnext-stytch-member/lambda/internal/service"
	"github.com/steverhoton/srnext-stytch-member/lambda/pkg/errors"
	"go.uber.org/zap"
)

// Handler handles Lambda requests
type Handler struct {
	memberService service.MemberService
	logger        *zap.Logger
	config        *config.Config
}

// NewHandler creates a new Lambda handler
func NewHandler(cfg *config.Config, logger *zap.Logger) (*Handler, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	memberService, err := service.NewMemberService(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create member service: %w", err)
	}

	return &Handler{
		memberService: memberService,
		logger:        logger,
		config:        cfg,
	}, nil
}

// HandleRequest handles ALB requests
func (h *Handler) HandleRequest(ctx context.Context, req models.ALBRequest) (models.ALBResponse, error) {
	h.logger.Info("Handling ALB request",
		zap.String("method", req.HTTPMethod),
		zap.String("path", req.Path))

	// Add CORS headers
	headers := map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type, Authorization",
	}

	// Handle OPTIONS request for CORS
	if req.HTTPMethod == http.MethodOptions {
		return models.ALBResponse{
			StatusCode: http.StatusOK,
			Headers:    headers,
		}, nil
	}

	// Route the request based on path and method
	response, err := h.routeRequest(ctx, req)
	if err != nil {
		h.logger.Error("Request failed", zap.Error(err))
		return h.errorResponse(err, headers), nil
	}

	// Marshal response
	body, err := json.Marshal(response)
	if err != nil {
		h.logger.Error("Failed to marshal response", zap.Error(err))
		return h.errorResponse(errors.ErrInternalServer, headers), nil
	}

	return models.ALBResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       string(body),
	}, nil
}

// routeRequest routes the request to the appropriate handler
func (h *Handler) routeRequest(ctx context.Context, req models.ALBRequest) (interface{}, error) {
	// Parse the path to determine the resource and operation
	pathParts := strings.Split(strings.Trim(req.Path, "/"), "/")

	// Expected paths:
	// POST   /members                    - Create member
	// GET    /members/{member_id}        - Get member
	// PUT    /members/{member_id}        - Update member
	// DELETE /members/{member_id}        - Delete member
	// POST   /members/search             - Search members

	if len(pathParts) == 0 || pathParts[0] != "members" {
		return nil, errors.NewAPIError(http.StatusNotFound, "NOT_FOUND", "Invalid path")
	}

	switch req.HTTPMethod {
	case http.MethodPost:
		if len(pathParts) == 2 && pathParts[1] == "search" {
			return h.handleSearchMembers(ctx, req)
		}
		return h.handleCreateMember(ctx, req)

	case http.MethodGet:
		if len(pathParts) < 2 {
			return nil, errors.NewAPIError(http.StatusBadRequest, "BAD_REQUEST", "Member ID required")
		}
		// Handle health check endpoint
		if pathParts[1] == "health" {
			return h.handleHealthCheck(ctx)
		}
		return h.handleGetMember(ctx, req, pathParts[1])

	case http.MethodPut:
		if len(pathParts) < 2 {
			return nil, errors.NewAPIError(http.StatusBadRequest, "BAD_REQUEST", "Member ID required")
		}
		return h.handleUpdateMember(ctx, req, pathParts[1])

	case http.MethodDelete:
		if len(pathParts) < 2 {
			return nil, errors.NewAPIError(http.StatusBadRequest, "BAD_REQUEST", "Member ID required")
		}
		return h.handleDeleteMember(ctx, req, pathParts[1])

	default:
		return nil, errors.NewAPIError(http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
	}
}

// handleHealthCheck handles health check requests
func (h *Handler) handleHealthCheck(ctx context.Context) (map[string]interface{}, error) {
	return map[string]interface{}{
		"status":  "healthy",
		"service": "stytch-member-management",
		"version": "1.0.0",
	}, nil
}

// handleCreateMember handles member creation
func (h *Handler) handleCreateMember(ctx context.Context, req models.ALBRequest) (*models.MemberResponse, error) {
	var createReq models.CreateMemberRequest
	if err := json.Unmarshal([]byte(req.Body), &createReq); err != nil {
		return nil, errors.NewAPIErrorWithDetails(http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", err.Error())
	}

	// Validate request
	if err := h.validateCreateRequest(&createReq); err != nil {
		return nil, err
	}

	return h.memberService.CreateMember(ctx, &createReq)
}

// handleGetMember handles member retrieval
func (h *Handler) handleGetMember(ctx context.Context, req models.ALBRequest, memberID string) (*models.MemberResponse, error) {
	// Get organization ID from query parameters or headers
	organizationID := h.getOrganizationID(req)
	if organizationID == "" {
		return nil, errors.NewAPIError(http.StatusBadRequest, "BAD_REQUEST", "Organization ID required")
	}

	getReq := &models.GetMemberRequest{
		OrganizationID: organizationID,
		MemberID:       memberID,
	}

	return h.memberService.GetMember(ctx, getReq)
}

// handleUpdateMember handles member updates
func (h *Handler) handleUpdateMember(ctx context.Context, req models.ALBRequest, memberID string) (*models.MemberResponse, error) {
	var updateReq models.UpdateMemberRequest
	if err := json.Unmarshal([]byte(req.Body), &updateReq); err != nil {
		return nil, errors.NewAPIErrorWithDetails(http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", err.Error())
	}

	// Set member ID from path
	updateReq.MemberID = memberID

	// Validate request
	if err := h.validateUpdateRequest(&updateReq); err != nil {
		return nil, err
	}

	return h.memberService.UpdateMember(ctx, &updateReq)
}

// handleDeleteMember handles member deletion
func (h *Handler) handleDeleteMember(ctx context.Context, req models.ALBRequest, memberID string) (*models.MemberResponse, error) {
	// Get organization ID from query parameters or headers
	organizationID := h.getOrganizationID(req)
	if organizationID == "" {
		return nil, errors.NewAPIError(http.StatusBadRequest, "BAD_REQUEST", "Organization ID required")
	}

	deleteReq := &models.DeleteMemberRequest{
		OrganizationID: organizationID,
		MemberID:       memberID,
	}

	return h.memberService.DeleteMember(ctx, deleteReq)
}

// handleSearchMembers handles member search
func (h *Handler) handleSearchMembers(ctx context.Context, req models.ALBRequest) (*models.MembersListResponse, error) {
	var searchReq models.SearchMembersRequest
	if err := json.Unmarshal([]byte(req.Body), &searchReq); err != nil {
		return nil, errors.NewAPIErrorWithDetails(http.StatusBadRequest, "BAD_REQUEST", "Invalid request body", err.Error())
	}

	// Validate request
	if len(searchReq.OrganizationIDs) == 0 {
		return nil, errors.NewAPIError(http.StatusBadRequest, "BAD_REQUEST", "At least one organization ID required")
	}

	return h.memberService.SearchMembers(ctx, &searchReq)
}

// getOrganizationID gets organization ID from request
func (h *Handler) getOrganizationID(req models.ALBRequest) string {
	// Check query parameters first
	if req.QueryStringParameters != nil {
		if orgID, ok := req.QueryStringParameters["organization_id"]; ok {
			return orgID
		}
	}

	// Check headers
	if req.Headers != nil {
		if orgID, ok := req.Headers["x-organization-id"]; ok {
			return orgID
		}
	}

	return ""
}

// validateCreateRequest validates create member request
func (h *Handler) validateCreateRequest(req *models.CreateMemberRequest) error {
	if req.OrganizationID == "" {
		return errors.NewAPIError(http.StatusBadRequest, "BAD_REQUEST", "Organization ID is required")
	}

	if req.EmailAddress == "" {
		return errors.NewAPIError(http.StatusBadRequest, "BAD_REQUEST", "Email address is required")
	}

	// Basic email validation
	if !strings.Contains(req.EmailAddress, "@") {
		return errors.NewAPIError(http.StatusBadRequest, "BAD_REQUEST", "Invalid email address")
	}

	return nil
}

// validateUpdateRequest validates update member request
func (h *Handler) validateUpdateRequest(req *models.UpdateMemberRequest) error {
	if req.OrganizationID == "" {
		return errors.NewAPIError(http.StatusBadRequest, "BAD_REQUEST", "Organization ID is required")
	}

	if req.MemberID == "" {
		return errors.NewAPIError(http.StatusBadRequest, "BAD_REQUEST", "Member ID is required")
	}

	// If email is provided, validate it
	if req.EmailAddress != "" && !strings.Contains(req.EmailAddress, "@") {
		return errors.NewAPIError(http.StatusBadRequest, "BAD_REQUEST", "Invalid email address")
	}

	return nil
}

// errorResponse creates an error response
func (h *Handler) errorResponse(err error, headers map[string]string) models.ALBResponse {
	var apiErr *errors.APIError
	var ok bool

	if apiErr, ok = err.(*errors.APIError); !ok {
		// Convert to API error if not already
		apiErr = errors.NewAPIErrorWithDetails(
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"An error occurred",
			err.Error(),
		)
	}

	errResp := errors.NewErrorResponse(apiErr)
	body, _ := json.Marshal(errResp)

	return models.ALBResponse{
		StatusCode: apiErr.StatusCode,
		Headers:    headers,
		Body:       string(body),
	}
}
