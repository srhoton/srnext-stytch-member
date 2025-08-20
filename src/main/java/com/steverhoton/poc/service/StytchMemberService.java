package com.steverhoton.poc.service;

import java.time.Instant;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.UUID;
import java.util.concurrent.ConcurrentHashMap;

import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.steverhoton.poc.exception.MemberNotFoundException;
import com.steverhoton.poc.exception.ValidationException;
import com.steverhoton.poc.model.MemberRequest;
import com.steverhoton.poc.stytch.mock.StytchMockModels;

/**
 * Service for managing Stytch members. This implementation uses mock models for demonstration. In
 * production, this would integrate with the actual Stytch SDK.
 */
public class StytchMemberService {

  private static final Logger logger = LoggerFactory.getLogger(StytchMemberService.class);
  private final String projectId;
  private final String secret;
  private final String organizationId;

  // Mock storage for demonstration
  private final Map<String, StytchMockModels.Member> memberStore = new ConcurrentHashMap<>();

  public StytchMemberService(String projectId, String secret, String organizationId) {
    if (StringUtils.isBlank(projectId)) {
      throw new IllegalArgumentException("Project ID cannot be blank");
    }
    if (StringUtils.isBlank(secret)) {
      throw new IllegalArgumentException("Secret cannot be blank");
    }
    if (StringUtils.isBlank(organizationId)) {
      throw new IllegalArgumentException("Organization ID cannot be blank");
    }

    this.projectId = projectId;
    this.secret = secret;
    this.organizationId = organizationId;
    logger.info("StytchMemberService initialized for organization: {}", organizationId);
  }

  public StytchMockModels.Member createMember(MemberRequest request) {
    validateMemberRequest(request);

    // Mock implementation
    StytchMockModels.Member member = new StytchMockModels.Member();
    member.setMemberId("member-" + UUID.randomUUID().toString());
    member.setOrganizationId(organizationId);
    member.setEmailAddress(request.getEmailAddress());
    member.setName(request.getName());
    member.setStatus("active");
    member.setTrustedMetadata(request.getTrustedMetadata());
    member.setUntrustedMetadata(request.getUntrustedMetadata());
    member.setIsBreakglass(request.getIsBreakglass());
    member.setMfaPhoneNumber(request.getMfaPhoneNumber());
    member.setMfaEnrolled(request.getMfaEnrolled());
    member.setCreatedAt(Instant.now());
    member.setUpdatedAt(Instant.now());

    memberStore.put(member.getMemberId(), member);
    logger.info("Successfully created member with ID: {}", member.getMemberId());

    return member;
  }

  public StytchMockModels.Member getMember(String memberId) {
    if (StringUtils.isBlank(memberId)) {
      throw new ValidationException("Member ID cannot be blank");
    }

    StytchMockModels.Member member = memberStore.get(memberId);
    if (member == null) {
      throw new MemberNotFoundException("Member not found with ID: " + memberId);
    }

    logger.info("Successfully retrieved member with ID: {}", memberId);
    return member;
  }

  public StytchMockModels.SearchResponse searchMembers(Map<String, String> queryParams) {
    StytchMockModels.SearchResponse response = new StytchMockModels.SearchResponse();
    List<StytchMockModels.Member> results = new ArrayList<>();

    if (queryParams != null) {
      String email = queryParams.get("email");
      if (StringUtils.isNotBlank(email)) {
        for (StytchMockModels.Member member : memberStore.values()) {
          if (email.equals(member.getEmailAddress())) {
            results.add(member);
          }
        }
      } else {
        results.addAll(memberStore.values());
      }

      String limit = queryParams.get("limit");
      if (StringUtils.isNotBlank(limit)) {
        int maxResults = Integer.parseInt(limit);
        if (results.size() > maxResults) {
          results = results.subList(0, maxResults);
        }
      }
    } else {
      results.addAll(memberStore.values());
    }

    response.setMembers(results);
    response.setStatusCode("200");
    response.setRequestId(UUID.randomUUID().toString());

    logger.info("Successfully searched members, found {} results", results.size());
    return response;
  }

  public StytchMockModels.Member updateMember(String memberId, MemberRequest request) {
    if (StringUtils.isBlank(memberId)) {
      throw new ValidationException("Member ID cannot be blank");
    }

    StytchMockModels.Member member = memberStore.get(memberId);
    if (member == null) {
      throw new MemberNotFoundException("Member not found with ID: " + memberId);
    }

    // Update fields if provided
    if (request.getEmailAddress() != null) {
      if (!isValidEmail(request.getEmailAddress())) {
        throw new ValidationException("Invalid email address format");
      }
      member.setEmailAddress(request.getEmailAddress());
    }
    if (request.getName() != null) {
      member.setName(request.getName());
    }
    if (request.getTrustedMetadata() != null) {
      member.setTrustedMetadata(request.getTrustedMetadata());
    }
    if (request.getUntrustedMetadata() != null) {
      member.setUntrustedMetadata(request.getUntrustedMetadata());
    }
    if (request.getIsBreakglass() != null) {
      member.setIsBreakglass(request.getIsBreakglass());
    }
    if (request.getMfaPhoneNumber() != null) {
      member.setMfaPhoneNumber(request.getMfaPhoneNumber());
    }
    if (request.getMfaEnrolled() != null) {
      member.setMfaEnrolled(request.getMfaEnrolled());
    }

    member.setUpdatedAt(Instant.now());
    logger.info("Successfully updated member with ID: {}", memberId);

    return member;
  }

  public void deleteMember(String memberId) {
    if (StringUtils.isBlank(memberId)) {
      throw new ValidationException("Member ID cannot be blank");
    }

    StytchMockModels.Member member = memberStore.remove(memberId);
    if (member == null) {
      throw new MemberNotFoundException("Member not found with ID: " + memberId);
    }

    logger.info("Successfully deleted member with ID: {}", memberId);
  }

  public StytchMockModels.Member reactivateMember(String memberId) {
    if (StringUtils.isBlank(memberId)) {
      throw new ValidationException("Member ID cannot be blank");
    }

    StytchMockModels.Member member = memberStore.get(memberId);
    if (member == null) {
      throw new MemberNotFoundException("Member not found with ID: " + memberId);
    }

    member.setStatus("active");
    member.setUpdatedAt(Instant.now());

    logger.info("Successfully reactivated member with ID: {}", memberId);
    return member;
  }

  private void validateMemberRequest(MemberRequest request) {
    if (request == null) {
      throw new ValidationException("Member request cannot be null");
    }
    if (StringUtils.isBlank(request.getEmailAddress())) {
      throw new ValidationException("Email address is required");
    }
    if (!isValidEmail(request.getEmailAddress())) {
      throw new ValidationException("Invalid email address format");
    }
  }

  private boolean isValidEmail(String email) {
    return email != null && email.matches("^[A-Za-z0-9+_.-]+@([A-Za-z0-9.-]+\\.[A-Za-z]{2,})$");
  }
}
