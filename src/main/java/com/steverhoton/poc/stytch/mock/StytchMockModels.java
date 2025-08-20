package com.steverhoton.poc.stytch.mock;

import java.time.Instant;
import java.util.List;
import java.util.Map;

/**
 * Mock models to demonstrate the Stytch Member API structure. In production, these would be
 * replaced with actual Stytch SDK models.
 */
public class StytchMockModels {

  public static class Member {
    private String memberId;
    private String organizationId;
    private String emailAddress;
    private String name;
    private String status;
    private Map<String, Object> trustedMetadata;
    private Map<String, Object> untrustedMetadata;
    private Boolean isBreakglass;
    private String mfaPhoneNumber;
    private Boolean mfaEnrolled;
    private Instant createdAt;
    private Instant updatedAt;

    public String getMemberId() {
      return memberId;
    }

    public void setMemberId(String memberId) {
      this.memberId = memberId;
    }

    public String getOrganizationId() {
      return organizationId;
    }

    public void setOrganizationId(String organizationId) {
      this.organizationId = organizationId;
    }

    public String getEmailAddress() {
      return emailAddress;
    }

    public void setEmailAddress(String emailAddress) {
      this.emailAddress = emailAddress;
    }

    public String getName() {
      return name;
    }

    public void setName(String name) {
      this.name = name;
    }

    public String getStatus() {
      return status;
    }

    public void setStatus(String status) {
      this.status = status;
    }

    public Map<String, Object> getTrustedMetadata() {
      return trustedMetadata;
    }

    public void setTrustedMetadata(Map<String, Object> trustedMetadata) {
      this.trustedMetadata = trustedMetadata;
    }

    public Map<String, Object> getUntrustedMetadata() {
      return untrustedMetadata;
    }

    public void setUntrustedMetadata(Map<String, Object> untrustedMetadata) {
      this.untrustedMetadata = untrustedMetadata;
    }

    public Boolean getIsBreakglass() {
      return isBreakglass;
    }

    public void setIsBreakglass(Boolean isBreakglass) {
      this.isBreakglass = isBreakglass;
    }

    public String getMfaPhoneNumber() {
      return mfaPhoneNumber;
    }

    public void setMfaPhoneNumber(String mfaPhoneNumber) {
      this.mfaPhoneNumber = mfaPhoneNumber;
    }

    public Boolean getMfaEnrolled() {
      return mfaEnrolled;
    }

    public void setMfaEnrolled(Boolean mfaEnrolled) {
      this.mfaEnrolled = mfaEnrolled;
    }

    public Instant getCreatedAt() {
      return createdAt;
    }

    public void setCreatedAt(Instant createdAt) {
      this.createdAt = createdAt;
    }

    public Instant getUpdatedAt() {
      return updatedAt;
    }

    public void setUpdatedAt(Instant updatedAt) {
      this.updatedAt = updatedAt;
    }
  }

  public static class SearchResponse {
    private List<Member> members;
    private String nextCursor;
    private String statusCode;
    private String requestId;

    public List<Member> getMembers() {
      return members;
    }

    public void setMembers(List<Member> members) {
      this.members = members;
    }

    public String getNextCursor() {
      return nextCursor;
    }

    public void setNextCursor(String nextCursor) {
      this.nextCursor = nextCursor;
    }

    public String getStatusCode() {
      return statusCode;
    }

    public void setStatusCode(String statusCode) {
      this.statusCode = statusCode;
    }

    public String getRequestId() {
      return requestId;
    }

    public void setRequestId(String requestId) {
      this.requestId = requestId;
    }
  }

  public static class StytchException extends Exception {
    private final int statusCode;

    public StytchException(String message, int statusCode) {
      super(message);
      this.statusCode = statusCode;
    }

    public int getStatusCode() {
      return statusCode;
    }
  }
}
