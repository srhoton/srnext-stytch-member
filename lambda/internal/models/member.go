package models

import "time"

// CreateMemberRequest represents the request to create a new member
type CreateMemberRequest struct {
	OrganizationID        string                 `json:"organization_id" validate:"required,uuid"`
	EmailAddress          string                 `json:"email_address" validate:"required,email"`
	Name                  string                 `json:"name,omitempty"`
	ExternalID            string                 `json:"external_id,omitempty"`
	Roles                 []string               `json:"roles,omitempty"`
	TrustedMetadata       map[string]interface{} `json:"trusted_metadata,omitempty"`
	UntrustedMetadata     map[string]interface{} `json:"untrusted_metadata,omitempty"`
	MFAPhoneNumber        string                 `json:"mfa_phone_number,omitempty"`
	CreateMemberAsPending bool                   `json:"create_member_as_pending,omitempty"`
	IsBreakglass          bool                   `json:"is_breakglass,omitempty"`
}

// UpdateMemberRequest represents the request to update a member
type UpdateMemberRequest struct {
	OrganizationID    string                 `json:"organization_id" validate:"required,uuid"`
	MemberID          string                 `json:"member_id" validate:"required,uuid"`
	EmailAddress      string                 `json:"email_address,omitempty" validate:"omitempty,email"`
	Name              string                 `json:"name,omitempty"`
	ExternalID        string                 `json:"external_id,omitempty"`
	Roles             []string               `json:"roles,omitempty"`
	TrustedMetadata   map[string]interface{} `json:"trusted_metadata,omitempty"`
	UntrustedMetadata map[string]interface{} `json:"untrusted_metadata,omitempty"`
	MFAPhoneNumber    string                 `json:"mfa_phone_number,omitempty"`
	IsBreakglass      bool                   `json:"is_breakglass,omitempty"`
	MFAEnrolled       bool                   `json:"mfa_enrolled,omitempty"`
	DefaultMFAMethod  string                 `json:"default_mfa_method,omitempty"`
}

// GetMemberRequest represents the request to get a member
type GetMemberRequest struct {
	OrganizationID string `json:"organization_id" validate:"required,uuid"`
	MemberID       string `json:"member_id,omitempty" validate:"omitempty,uuid"`
	EmailAddress   string `json:"email_address,omitempty" validate:"omitempty,email"`
}

// DeleteMemberRequest represents the request to delete a member
type DeleteMemberRequest struct {
	OrganizationID string `json:"organization_id" validate:"required,uuid"`
	MemberID       string `json:"member_id" validate:"required,uuid"`
}

// SearchMembersRequest represents the request to search members
type SearchMembersRequest struct {
	OrganizationIDs []string          `json:"organization_ids" validate:"required,min=1,dive,uuid"`
	Cursor          string            `json:"cursor,omitempty"`
	Limit           int               `json:"limit,omitempty" validate:"omitempty,min=1,max=1000"`
	Query           map[string]string `json:"query,omitempty"`
}

// Member represents a Stytch member
type Member struct {
	MemberID          string                 `json:"member_id"`
	OrganizationID    string                 `json:"organization_id"`
	EmailAddress      string                 `json:"email_address"`
	Status            string                 `json:"status"`
	Name              string                 `json:"name,omitempty"`
	ExternalID        string                 `json:"external_id,omitempty"`
	Roles             []string               `json:"roles,omitempty"`
	TrustedMetadata   map[string]interface{} `json:"trusted_metadata,omitempty"`
	UntrustedMetadata map[string]interface{} `json:"untrusted_metadata,omitempty"`
	MFAPhoneNumber    string                 `json:"mfa_phone_number,omitempty"`
	MFAEnrolled       bool                   `json:"mfa_enrolled"`
	DefaultMFAMethod  string                 `json:"default_mfa_method,omitempty"`
	IsBreakglass      bool                   `json:"is_breakglass"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

// MemberResponse represents the response for member operations
type MemberResponse struct {
	StatusCode int    `json:"status_code"`
	RequestID  string `json:"request_id"`
	Member     Member `json:"member,omitempty"`
}

// MembersListResponse represents the response for search operations
type MembersListResponse struct {
	StatusCode int      `json:"status_code"`
	RequestID  string   `json:"request_id"`
	Members    []Member `json:"members"`
	Cursor     string   `json:"cursor,omitempty"`
}
