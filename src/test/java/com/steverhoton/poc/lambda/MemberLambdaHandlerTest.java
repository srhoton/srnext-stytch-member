package com.steverhoton.poc.lambda;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.*;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import com.amazonaws.services.lambda.runtime.Context;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.steverhoton.poc.exception.MemberNotFoundException;
import com.steverhoton.poc.exception.StytchApiException;
import com.steverhoton.poc.exception.ValidationException;
import com.steverhoton.poc.model.ApiGatewayRequest;
import com.steverhoton.poc.model.ApiGatewayResponse;
import com.steverhoton.poc.model.MemberRequest;
import com.steverhoton.poc.service.StytchMemberService;
import com.steverhoton.poc.stytch.mock.StytchMockModels;

import uk.org.webcompere.systemstubs.environment.EnvironmentVariables;
import uk.org.webcompere.systemstubs.jupiter.SystemStub;
import uk.org.webcompere.systemstubs.jupiter.SystemStubsExtension;

@ExtendWith({MockitoExtension.class, SystemStubsExtension.class})
class MemberLambdaHandlerTest {

  @Mock private StytchMemberService memberService;

  @Mock private Context context;

  @SystemStub private EnvironmentVariables environmentVariables;

  private MemberLambdaHandler handler;
  private ObjectMapper objectMapper = new ObjectMapper();

  @BeforeEach
  void setUp() {
    handler = new MemberLambdaHandler(memberService);
  }

  @Test
  void constructor_throwsException_whenProjectIdNotSet() throws Exception {
    environmentVariables.set("STYTCH_SECRET", "secret");
    environmentVariables.set("STYTCH_ORGANIZATION_ID", "org");

    assertThatThrownBy(() -> new MemberLambdaHandler())
        .isInstanceOf(IllegalStateException.class)
        .hasMessage("STYTCH_PROJECT_ID environment variable is not set");
  }

  @Test
  void constructor_throwsException_whenSecretNotSet() throws Exception {
    environmentVariables.set("STYTCH_PROJECT_ID", "project");
    environmentVariables.set("STYTCH_ORGANIZATION_ID", "org");

    assertThatThrownBy(() -> new MemberLambdaHandler())
        .isInstanceOf(IllegalStateException.class)
        .hasMessage("STYTCH_SECRET environment variable is not set");
  }

  @Test
  void handleRequest_returnsError_whenMethodIsNull() {
    ApiGatewayRequest request = new ApiGatewayRequest();
    request.setPath("/members");

    ApiGatewayResponse response = handler.handleRequest(request, context);

    assertThat(response.getStatusCode()).isEqualTo(400);
    assertThat(response.getBody()).contains("Missing method or path");
  }

  @Test
  void handleRequest_returnsError_whenPathIsNull() {
    ApiGatewayRequest request = new ApiGatewayRequest();
    request.setHttpMethod("GET");

    ApiGatewayResponse response = handler.handleRequest(request, context);

    assertThat(response.getStatusCode()).isEqualTo(400);
    assertThat(response.getBody()).contains("Missing method or path");
  }

  @Test
  void handleRequest_returnsNotFound_forUnknownPath() {
    ApiGatewayRequest request = new ApiGatewayRequest();
    request.setHttpMethod("GET");
    request.setPath("/unknown");

    ApiGatewayResponse response = handler.handleRequest(request, context);

    assertThat(response.getStatusCode()).isEqualTo(404);
    assertThat(response.getBody()).contains("Endpoint not found");
  }

  @Test
  void handleRequest_createMember_success() throws Exception {
    ApiGatewayRequest request = new ApiGatewayRequest();
    request.setHttpMethod("POST");
    request.setPath("/members");

    MemberRequest memberRequest = new MemberRequest();
    memberRequest.setEmailAddress("test@example.com");
    memberRequest.setName("Test User");
    request.setBody(objectMapper.writeValueAsString(memberRequest));

    StytchMockModels.Member member = new StytchMockModels.Member();
    member.setMemberId("member-123");
    member.setEmailAddress("test@example.com");

    when(memberService.createMember(any(MemberRequest.class))).thenReturn(member);

    ApiGatewayResponse response = handler.handleRequest(request, context);

    assertThat(response.getStatusCode()).isEqualTo(201);
    assertThat(response.getBody()).contains("member-123");
    assertThat(response.getBody()).contains("test@example.com");
  }

  @Test
  void handleRequest_createMember_returnsError_whenBodyIsEmpty() {
    ApiGatewayRequest request = new ApiGatewayRequest();
    request.setHttpMethod("POST");
    request.setPath("/members");
    request.setBody("");

    ApiGatewayResponse response = handler.handleRequest(request, context);

    assertThat(response.getStatusCode()).isEqualTo(400);
    assertThat(response.getBody()).contains("Request body is required");
  }

  @Test
  void handleRequest_getMember_success() throws Exception {
    ApiGatewayRequest request = new ApiGatewayRequest();
    request.setHttpMethod("GET");
    request.setPath("/members/member-123");

    StytchMockModels.Member member = new StytchMockModels.Member();
    member.setMemberId("member-123");
    member.setEmailAddress("test@example.com");

    when(memberService.getMember("member-123")).thenReturn(member);

    ApiGatewayResponse response = handler.handleRequest(request, context);

    assertThat(response.getStatusCode()).isEqualTo(200);
    assertThat(response.getBody()).contains("member-123");
    assertThat(response.getBody()).contains("test@example.com");
  }

  @Test
  void handleRequest_getMember_returnsNotFound() {
    ApiGatewayRequest request = new ApiGatewayRequest();
    request.setHttpMethod("GET");
    request.setPath("/members/nonexistent");

    when(memberService.getMember("nonexistent"))
        .thenThrow(new MemberNotFoundException("Member not found"));

    ApiGatewayResponse response = handler.handleRequest(request, context);

    assertThat(response.getStatusCode()).isEqualTo(404);
    assertThat(response.getBody()).contains("Member not found");
  }

  @Test
  void handleRequest_searchMembers_success() throws Exception {
    ApiGatewayRequest request = new ApiGatewayRequest();
    request.setHttpMethod("GET");
    request.setPath("/members");

    Map<String, String> queryParams = new HashMap<>();
    queryParams.put("email", "test@example.com");
    request.setQueryStringParameters(queryParams);

    StytchMockModels.SearchResponse searchResponse = new StytchMockModels.SearchResponse();
    searchResponse.setMembers(List.of());

    when(memberService.searchMembers(queryParams)).thenReturn(searchResponse);

    ApiGatewayResponse response = handler.handleRequest(request, context);

    assertThat(response.getStatusCode()).isEqualTo(200);
    verify(memberService).searchMembers(queryParams);
  }

  @Test
  void handleRequest_updateMember_success() throws Exception {
    ApiGatewayRequest request = new ApiGatewayRequest();
    request.setHttpMethod("PUT");
    request.setPath("/members/member-123");

    MemberRequest memberRequest = new MemberRequest();
    memberRequest.setName("Updated Name");
    request.setBody(objectMapper.writeValueAsString(memberRequest));

    StytchMockModels.Member member = new StytchMockModels.Member();
    member.setMemberId("member-123");
    member.setName("Updated Name");

    when(memberService.updateMember(eq("member-123"), any(MemberRequest.class))).thenReturn(member);

    ApiGatewayResponse response = handler.handleRequest(request, context);

    assertThat(response.getStatusCode()).isEqualTo(200);
    assertThat(response.getBody()).contains("Updated Name");
  }

  @Test
  void handleRequest_deleteMember_success() {
    ApiGatewayRequest request = new ApiGatewayRequest();
    request.setHttpMethod("DELETE");
    request.setPath("/members/member-123");

    doNothing().when(memberService).deleteMember("member-123");

    ApiGatewayResponse response = handler.handleRequest(request, context);

    assertThat(response.getStatusCode()).isEqualTo(204);
    assertThat(response.getBody()).isNull();
    verify(memberService).deleteMember("member-123");
  }

  @Test
  void handleRequest_reactivateMember_success() throws Exception {
    ApiGatewayRequest request = new ApiGatewayRequest();
    request.setHttpMethod("POST");
    request.setPath("/members/member-123/reactivate");

    StytchMockModels.Member member = new StytchMockModels.Member();
    member.setMemberId("member-123");
    member.setEmailAddress("test@example.com");

    when(memberService.reactivateMember("member-123")).thenReturn(member);

    ApiGatewayResponse response = handler.handleRequest(request, context);

    assertThat(response.getStatusCode()).isEqualTo(200);
    assertThat(response.getBody()).contains("member-123");
  }

  @Test
  void handleRequest_returnsMethodNotAllowed_forInvalidMethod() {
    ApiGatewayRequest request = new ApiGatewayRequest();
    request.setHttpMethod("PATCH");
    request.setPath("/members");

    ApiGatewayResponse response = handler.handleRequest(request, context);

    assertThat(response.getStatusCode()).isEqualTo(405);
    assertThat(response.getBody()).contains("Method PATCH not allowed");
  }

  @Test
  void handleRequest_handlesValidationException() {
    ApiGatewayRequest request = new ApiGatewayRequest();
    request.setHttpMethod("POST");
    request.setPath("/members");
    request.setBody("{\"email_address\": \"invalid\"}");

    when(memberService.createMember(any()))
        .thenThrow(new ValidationException("Invalid email format"));

    ApiGatewayResponse response = handler.handleRequest(request, context);

    assertThat(response.getStatusCode()).isEqualTo(400);
    assertThat(response.getBody()).contains("Invalid email format");
  }

  @Test
  void handleRequest_handlesStytchApiException() {
    ApiGatewayRequest request = new ApiGatewayRequest();
    request.setHttpMethod("GET");
    request.setPath("/members/member-123");

    when(memberService.getMember(anyString())).thenThrow(new StytchApiException("API Error", 503));

    ApiGatewayResponse response = handler.handleRequest(request, context);

    assertThat(response.getStatusCode()).isEqualTo(503);
    assertThat(response.getBody()).contains("API Error");
  }

  @Test
  void handleRequest_handlesUnexpectedException() {
    ApiGatewayRequest request = new ApiGatewayRequest();
    request.setHttpMethod("GET");
    request.setPath("/members/member-123");

    when(memberService.getMember(anyString())).thenThrow(new RuntimeException("Unexpected error"));

    ApiGatewayResponse response = handler.handleRequest(request, context);

    assertThat(response.getStatusCode()).isEqualTo(500);
    assertThat(response.getBody()).contains("An unexpected error occurred");
  }
}
