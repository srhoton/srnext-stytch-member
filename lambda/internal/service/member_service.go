package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/steverhoton/srnext-stytch-member/lambda/internal/config"
	"github.com/steverhoton/srnext-stytch-member/lambda/internal/models"
	"github.com/steverhoton/srnext-stytch-member/lambda/pkg/errors"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/b2bstytchapi"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/organizations"
	"github.com/stytchauth/stytch-go/v12/stytch/b2b/organizations/members"
	"go.uber.org/zap"
)

// MemberService handles Stytch member operations
type MemberService interface {
	CreateMember(ctx context.Context, req *models.CreateMemberRequest) (*models.MemberResponse, error)
	GetMember(ctx context.Context, req *models.GetMemberRequest) (*models.MemberResponse, error)
	UpdateMember(ctx context.Context, req *models.UpdateMemberRequest) (*models.MemberResponse, error)
	DeleteMember(ctx context.Context, req *models.DeleteMemberRequest) (*models.MemberResponse, error)
	SearchMembers(ctx context.Context, req *models.SearchMembersRequest) (*models.MembersListResponse, error)
}

// stytchMemberService implements MemberService
type stytchMemberService struct {
	client *b2bstytchapi.API
	logger *zap.Logger
	config *config.Config
}

// NewMemberService creates a new member service
func NewMemberService(cfg *config.Config, logger *zap.Logger) (MemberService, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}
	
	// Create Stytch client
	// Note: The SDK automatically uses test.stytch.com for test credentials
	client, err := b2bstytchapi.NewClient(cfg.StytchProjectID, cfg.StytchSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stytch client: %w", err)
	}
	
	return &stytchMemberService{
		client: client,
		logger: logger,
		config: cfg,
	}, nil
}

// CreateMember creates a new member in Stytch
func (s *stytchMemberService) CreateMember(ctx context.Context, req *models.CreateMemberRequest) (*models.MemberResponse, error) {
	s.logger.Info("Creating member",
		zap.String("organization_id", req.OrganizationID),
		zap.String("email", req.EmailAddress))
	
	createReq := &members.CreateParams{
		OrganizationID:        req.OrganizationID,
		EmailAddress:          req.EmailAddress,
		Name:                  req.Name,
		CreateMemberAsPending: req.CreateMemberAsPending,
		IsBreakglass:          req.IsBreakglass,
	}
	
	// Add optional fields
	if len(req.Roles) > 0 {
		// Convert string roles to proper type if needed
		createReq.Roles = req.Roles
	}
	
	if req.TrustedMetadata != nil {
		createReq.TrustedMetadata = req.TrustedMetadata
	}
	
	if req.UntrustedMetadata != nil {
		createReq.UntrustedMetadata = req.UntrustedMetadata
	}
	
	if req.MFAPhoneNumber != "" {
		createReq.MFAPhoneNumber = req.MFAPhoneNumber
	}
	
	resp, err := s.client.Organizations.Members.Create(ctx, createReq)
	if err != nil {
		s.logger.Error("Failed to create member",
			zap.String("organization_id", req.OrganizationID),
			zap.String("email", req.EmailAddress),
			zap.Error(err))
		return nil, s.handleStytchError(err)
	}
	
	return s.mapCreateResponse(resp), nil
}

// GetMember retrieves a member from Stytch
func (s *stytchMemberService) GetMember(ctx context.Context, req *models.GetMemberRequest) (*models.MemberResponse, error) {
	s.logger.Info("Getting member",
		zap.String("organization_id", req.OrganizationID),
		zap.String("member_id", req.MemberID),
		zap.String("email", req.EmailAddress))
	
	getReq := &members.GetParams{
		OrganizationID: req.OrganizationID,
	}
	
	// Stytch expects these as strings not pointers
	if req.MemberID != "" {
		getReq.MemberID = req.MemberID
	} else if req.EmailAddress != "" {
		getReq.EmailAddress = req.EmailAddress
	} else {
		return nil, errors.NewAPIError(http.StatusBadRequest, "BAD_REQUEST", "Either member_id or email_address is required")
	}
	
	resp, err := s.client.Organizations.Members.Get(ctx, getReq)
	if err != nil {
		s.logger.Error("Failed to get member",
			zap.String("organization_id", req.OrganizationID),
			zap.Error(err))
		return nil, s.handleStytchError(err)
	}
	
	return s.mapGetResponse(resp), nil
}

// UpdateMember updates a member in Stytch
func (s *stytchMemberService) UpdateMember(ctx context.Context, req *models.UpdateMemberRequest) (*models.MemberResponse, error) {
	s.logger.Info("Updating member",
		zap.String("organization_id", req.OrganizationID),
		zap.String("member_id", req.MemberID))
	
	updateReq := &members.UpdateParams{
		OrganizationID: req.OrganizationID,
		MemberID:       req.MemberID,
	}
	
	// Only set fields that are provided
	if req.EmailAddress != "" {
		updateReq.EmailAddress = req.EmailAddress
	}
	
	if req.Name != "" {
		updateReq.Name = req.Name
	}
	
	if len(req.Roles) > 0 {
		updateReq.Roles = req.Roles
	}
	
	if req.TrustedMetadata != nil {
		updateReq.TrustedMetadata = req.TrustedMetadata
	}
	
	if req.UntrustedMetadata != nil {
		updateReq.UntrustedMetadata = req.UntrustedMetadata
	}
	
	if req.MFAPhoneNumber != "" {
		updateReq.MFAPhoneNumber = req.MFAPhoneNumber
	}
	
	// Handle boolean fields explicitly
	updateReq.IsBreakglass = req.IsBreakglass
	updateReq.MFAEnrolled = req.MFAEnrolled
	
	if req.DefaultMFAMethod != "" {
		updateReq.DefaultMFAMethod = req.DefaultMFAMethod
	}
	
	resp, err := s.client.Organizations.Members.Update(ctx, updateReq)
	if err != nil {
		s.logger.Error("Failed to update member",
			zap.String("organization_id", req.OrganizationID),
			zap.String("member_id", req.MemberID),
			zap.Error(err))
		return nil, s.handleStytchError(err)
	}
	
	return s.mapUpdateResponse(resp), nil
}

// DeleteMember deletes a member from Stytch
func (s *stytchMemberService) DeleteMember(ctx context.Context, req *models.DeleteMemberRequest) (*models.MemberResponse, error) {
	s.logger.Info("Deleting member",
		zap.String("organization_id", req.OrganizationID),
		zap.String("member_id", req.MemberID))
	
	deleteReq := &members.DeleteParams{
		OrganizationID: req.OrganizationID,
		MemberID:       req.MemberID,
	}
	
	resp, err := s.client.Organizations.Members.Delete(ctx, deleteReq)
	if err != nil {
		s.logger.Error("Failed to delete member",
			zap.String("organization_id", req.OrganizationID),
			zap.String("member_id", req.MemberID),
			zap.Error(err))
		return nil, s.handleStytchError(err)
	}
	
	return &models.MemberResponse{
		StatusCode: int(resp.StatusCode),
		RequestID:  resp.RequestID,
		Member:     models.Member{MemberID: resp.MemberID}, // Stytch returns only the deleted member ID
	}, nil
}

// SearchMembers searches for members in Stytch
func (s *stytchMemberService) SearchMembers(ctx context.Context, req *models.SearchMembersRequest) (*models.MembersListResponse, error) {
	s.logger.Info("Searching members",
		zap.Strings("organization_ids", req.OrganizationIDs),
		zap.String("cursor", req.Cursor))
	
	searchReq := &members.SearchParams{
		OrganizationIds: req.OrganizationIDs, // Note: Ids not IDs
		Cursor:          req.Cursor,
	}
	
	if req.Limit > 0 {
		limit := uint32(req.Limit)
		searchReq.Limit = limit
	}
	
	// Note: Query structure may vary by SDK version
	// Simplified for now
	
	resp, err := s.client.Organizations.Members.Search(ctx, searchReq)
	if err != nil {
		s.logger.Error("Failed to search members",
			zap.Strings("organization_ids", req.OrganizationIDs),
			zap.Error(err))
		return nil, s.handleStytchError(err)
	}
	
	members := make([]models.Member, 0, len(resp.Members))
	for _, m := range resp.Members {
		members = append(members, s.mapMember(m))
	}
	
	return &models.MembersListResponse{
		StatusCode: int(resp.StatusCode),
		RequestID:  resp.RequestID,
		Members:    members,
		Cursor:     resp.ResultsMetadata.NextCursor,
	}, nil
}

// mapMember maps a Stytch member to our internal model
func (s *stytchMemberService) mapMember(m organizations.Member) models.Member {
	member := models.Member{
		MemberID:          m.MemberID,
		OrganizationID:    m.OrganizationID,
		EmailAddress:      m.EmailAddress,
		Status:            m.Status,
		Name:              m.Name,
		TrustedMetadata:   m.TrustedMetadata,
		UntrustedMetadata: m.UntrustedMetadata,
		MFAPhoneNumber:    m.MFAPhoneNumber,
		MFAEnrolled:       m.MFAEnrolled,
		IsBreakglass:      m.IsBreakglass,
		// Note: CreatedAt and UpdatedAt might not be available in this SDK version
	}
	
	// Map roles (convert from MemberRole to string if needed)
	if m.Roles != nil {
		roles := make([]string, 0, len(m.Roles))
		for _, role := range m.Roles {
			// Assuming MemberRole has a field we can use
			// This may need adjustment based on actual type
			roles = append(roles, string(role.RoleID))
		}
		member.Roles = roles
	}
	
	// Map DefaultMFAMethod if available
	if m.DefaultMFAMethod != "" {
		member.DefaultMFAMethod = m.DefaultMFAMethod
	}
	
	// ExternalID might not exist in this version
	// member.ExternalID = m.ExternalID
	
	return member
}

// Helper methods for mapping responses
func (s *stytchMemberService) mapCreateResponse(resp *members.CreateResponse) *models.MemberResponse {
	return &models.MemberResponse{
		StatusCode: int(resp.StatusCode),
		RequestID:  resp.RequestID,
		Member:     s.mapMember(resp.Member),
	}
}

func (s *stytchMemberService) mapGetResponse(resp *members.GetResponse) *models.MemberResponse {
	return &models.MemberResponse{
		StatusCode: int(resp.StatusCode),
		RequestID:  resp.RequestID,
		Member:     s.mapMember(resp.Member),
	}
}

func (s *stytchMemberService) mapUpdateResponse(resp *members.UpdateResponse) *models.MemberResponse {
	return &models.MemberResponse{
		StatusCode: int(resp.StatusCode),
		RequestID:  resp.RequestID,
		Member:     s.mapMember(resp.Member),
	}
}

// handleStytchError handles Stytch API errors
func (s *stytchMemberService) handleStytchError(err error) error {
	// For now, return a generic error wrapper
	// The Stytch SDK may handle errors differently in v12
	// We'll check the error message to determine the type
	if err == nil {
		return nil
	}
	
	errMsg := err.Error()
	
	// Try to determine error type from message
	if strings.Contains(errMsg, "400") || strings.Contains(errMsg, "bad request") {
		return errors.NewAPIErrorWithDetails(http.StatusBadRequest, "BAD_REQUEST", "Invalid request", errMsg)
	}
	if strings.Contains(errMsg, "401") || strings.Contains(errMsg, "unauthorized") {
		return errors.NewAPIErrorWithDetails(http.StatusUnauthorized, "UNAUTHORIZED", "Authentication failed", errMsg)
	}
	if strings.Contains(errMsg, "403") || strings.Contains(errMsg, "forbidden") {
		return errors.NewAPIErrorWithDetails(http.StatusForbidden, "FORBIDDEN", "Access denied", errMsg)
	}
	if strings.Contains(errMsg, "404") || strings.Contains(errMsg, "not found") {
		return errors.NewAPIErrorWithDetails(http.StatusNotFound, "NOT_FOUND", "Member not found", errMsg)
	}
	if strings.Contains(errMsg, "409") || strings.Contains(errMsg, "conflict") {
		return errors.NewAPIErrorWithDetails(http.StatusConflict, "CONFLICT", "Member already exists", errMsg)
	}
	if strings.Contains(errMsg, "429") || strings.Contains(errMsg, "rate") {
		return errors.NewAPIErrorWithDetails(http.StatusTooManyRequests, "RATE_LIMITED", "Too many requests", errMsg)
	}
	
	// Generic error
	return errors.NewAPIErrorWithDetails(http.StatusInternalServerError, "INTERNAL_ERROR", "An error occurred", errMsg)
}