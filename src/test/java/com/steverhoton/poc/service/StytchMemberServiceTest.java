package com.steverhoton.poc.service;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

import java.util.HashMap;
import java.util.Map;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import com.steverhoton.poc.exception.MemberNotFoundException;
import com.steverhoton.poc.exception.ValidationException;
import com.steverhoton.poc.model.MemberRequest;
import com.steverhoton.poc.stytch.mock.StytchMockModels;

class StytchMemberServiceTest {

  private StytchMemberService service;

  @BeforeEach
  void setUp() {
    service = new StytchMemberService("test-project-id", "test-secret", "test-org-id");
  }

  @Test
  void constructor_throwsException_whenProjectIdIsBlank() {
    assertThatThrownBy(() -> new StytchMemberService("", "secret", "org"))
        .isInstanceOf(IllegalArgumentException.class)
        .hasMessage("Project ID cannot be blank");
  }

  @Test
  void constructor_throwsException_whenSecretIsBlank() {
    assertThatThrownBy(() -> new StytchMemberService("project", "", "org"))
        .isInstanceOf(IllegalArgumentException.class)
        .hasMessage("Secret cannot be blank");
  }

  @Test
  void constructor_throwsException_whenOrganizationIdIsBlank() {
    assertThatThrownBy(() -> new StytchMemberService("project", "secret", ""))
        .isInstanceOf(IllegalArgumentException.class)
        .hasMessage("Organization ID cannot be blank");
  }

  @Test
  void createMember_throwsValidationException_whenRequestIsNull() {
    assertThatThrownBy(() -> service.createMember(null))
        .isInstanceOf(ValidationException.class)
        .hasMessage("Member request cannot be null");
  }

  @Test
  void createMember_throwsValidationException_whenEmailIsBlank() {
    MemberRequest request = new MemberRequest();
    request.setEmailAddress("");

    assertThatThrownBy(() -> service.createMember(request))
        .isInstanceOf(ValidationException.class)
        .hasMessage("Email address is required");
  }

  @Test
  void createMember_throwsValidationException_whenEmailIsInvalid() {
    MemberRequest request = new MemberRequest();
    request.setEmailAddress("invalid-email");

    assertThatThrownBy(() -> service.createMember(request))
        .isInstanceOf(ValidationException.class)
        .hasMessage("Invalid email address format");
  }

  @Test
  void getMember_throwsValidationException_whenMemberIdIsBlank() {
    assertThatThrownBy(() -> service.getMember(""))
        .isInstanceOf(ValidationException.class)
        .hasMessage("Member ID cannot be blank");
  }

  @Test
  void updateMember_throwsValidationException_whenMemberIdIsBlank() {
    MemberRequest request = new MemberRequest();
    request.setEmailAddress("test@example.com");

    assertThatThrownBy(() -> service.updateMember("", request))
        .isInstanceOf(ValidationException.class)
        .hasMessage("Member ID cannot be blank");
  }

  @Test
  void deleteMember_throwsValidationException_whenMemberIdIsBlank() {
    assertThatThrownBy(() -> service.deleteMember(""))
        .isInstanceOf(ValidationException.class)
        .hasMessage("Member ID cannot be blank");
  }

  @Test
  void reactivateMember_throwsValidationException_whenMemberIdIsBlank() {
    assertThatThrownBy(() -> service.reactivateMember(""))
        .isInstanceOf(ValidationException.class)
        .hasMessage("Member ID cannot be blank");
  }

  @Test
  void createMember_validatesEmailFormat() {
    MemberRequest request = new MemberRequest();

    request.setEmailAddress("notanemail");
    assertThatThrownBy(() -> service.createMember(request))
        .isInstanceOf(ValidationException.class)
        .hasMessage("Invalid email address format");

    request.setEmailAddress("@example.com");
    assertThatThrownBy(() -> service.createMember(request))
        .isInstanceOf(ValidationException.class)
        .hasMessage("Invalid email address format");

    request.setEmailAddress("user@");
    assertThatThrownBy(() -> service.createMember(request))
        .isInstanceOf(ValidationException.class)
        .hasMessage("Invalid email address format");

    request.setEmailAddress("user@domain");
    assertThatThrownBy(() -> service.createMember(request))
        .isInstanceOf(ValidationException.class)
        .hasMessage("Invalid email address format");
  }

  @Test
  void createMember_success() {
    MemberRequest request = new MemberRequest();
    request.setEmailAddress("test@example.com");
    request.setName("Test User");

    StytchMockModels.Member member = service.createMember(request);

    assertThat(member).isNotNull();
    assertThat(member.getMemberId()).startsWith("member-");
    assertThat(member.getEmailAddress()).isEqualTo("test@example.com");
    assertThat(member.getName()).isEqualTo("Test User");
    assertThat(member.getStatus()).isEqualTo("active");
  }

  @Test
  void getMember_success() {
    // First create a member
    MemberRequest request = new MemberRequest();
    request.setEmailAddress("test@example.com");
    StytchMockModels.Member created = service.createMember(request);

    // Then retrieve it
    StytchMockModels.Member retrieved = service.getMember(created.getMemberId());

    assertThat(retrieved).isNotNull();
    assertThat(retrieved.getMemberId()).isEqualTo(created.getMemberId());
    assertThat(retrieved.getEmailAddress()).isEqualTo("test@example.com");
  }

  @Test
  void getMember_throwsNotFoundException_whenNotExists() {
    assertThatThrownBy(() -> service.getMember("nonexistent"))
        .isInstanceOf(MemberNotFoundException.class)
        .hasMessage("Member not found with ID: nonexistent");
  }

  @Test
  void searchMembers_success() {
    // Create test members
    MemberRequest request1 = new MemberRequest();
    request1.setEmailAddress("test1@example.com");
    service.createMember(request1);

    MemberRequest request2 = new MemberRequest();
    request2.setEmailAddress("test2@example.com");
    service.createMember(request2);

    // Search without filters
    StytchMockModels.SearchResponse response = service.searchMembers(null);
    assertThat(response.getMembers()).hasSize(2);

    // Search with email filter
    Map<String, String> queryParams = new HashMap<>();
    queryParams.put("email", "test1@example.com");
    response = service.searchMembers(queryParams);
    assertThat(response.getMembers()).hasSize(1);
    assertThat(response.getMembers().get(0).getEmailAddress()).isEqualTo("test1@example.com");
  }

  @Test
  void updateMember_success() {
    // Create a member
    MemberRequest createRequest = new MemberRequest();
    createRequest.setEmailAddress("test@example.com");
    createRequest.setName("Original Name");
    StytchMockModels.Member created = service.createMember(createRequest);

    // Update it
    MemberRequest updateRequest = new MemberRequest();
    updateRequest.setName("Updated Name");
    StytchMockModels.Member updated = service.updateMember(created.getMemberId(), updateRequest);

    assertThat(updated.getName()).isEqualTo("Updated Name");
    assertThat(updated.getEmailAddress()).isEqualTo("test@example.com");
  }

  @Test
  void updateMember_throwsNotFoundException_whenNotExists() {
    MemberRequest request = new MemberRequest();
    request.setName("Test");

    assertThatThrownBy(() -> service.updateMember("nonexistent", request))
        .isInstanceOf(MemberNotFoundException.class)
        .hasMessage("Member not found with ID: nonexistent");
  }

  @Test
  void deleteMember_success() {
    // Create a member
    MemberRequest request = new MemberRequest();
    request.setEmailAddress("test@example.com");
    StytchMockModels.Member created = service.createMember(request);

    // Delete it
    service.deleteMember(created.getMemberId());

    // Verify it's gone
    assertThatThrownBy(() -> service.getMember(created.getMemberId()))
        .isInstanceOf(MemberNotFoundException.class);
  }

  @Test
  void deleteMember_throwsNotFoundException_whenNotExists() {
    assertThatThrownBy(() -> service.deleteMember("nonexistent"))
        .isInstanceOf(MemberNotFoundException.class)
        .hasMessage("Member not found with ID: nonexistent");
  }

  @Test
  void reactivateMember_success() {
    // Create a member
    MemberRequest request = new MemberRequest();
    request.setEmailAddress("test@example.com");
    StytchMockModels.Member created = service.createMember(request);

    // Reactivate it
    StytchMockModels.Member reactivated = service.reactivateMember(created.getMemberId());

    assertThat(reactivated.getStatus()).isEqualTo("active");
  }

  @Test
  void reactivateMember_throwsNotFoundException_whenNotExists() {
    assertThatThrownBy(() -> service.reactivateMember("nonexistent"))
        .isInstanceOf(MemberNotFoundException.class)
        .hasMessage("Member not found with ID: nonexistent");
  }
}
