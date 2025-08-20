package com.steverhoton.poc.model;

import java.util.Map;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public class MemberRequest {

  @JsonProperty("email_address")
  private String emailAddress;

  @JsonProperty("name")
  private String name;

  @JsonProperty("trusted_metadata")
  private Map<String, Object> trustedMetadata;

  @JsonProperty("untrusted_metadata")
  private Map<String, Object> untrustedMetadata;

  @JsonProperty("create_member_as_pending")
  private Boolean createMemberAsPending;

  @JsonProperty("is_breakglass")
  private Boolean isBreakglass;

  @JsonProperty("mfa_phone_number")
  private String mfaPhoneNumber;

  @JsonProperty("mfa_enrolled")
  private Boolean mfaEnrolled;

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

  public Boolean getCreateMemberAsPending() {
    return createMemberAsPending;
  }

  public void setCreateMemberAsPending(Boolean createMemberAsPending) {
    this.createMemberAsPending = createMemberAsPending;
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
}
