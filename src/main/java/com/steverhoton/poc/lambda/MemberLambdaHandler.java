package com.steverhoton.poc.lambda;

import java.util.Map;

import org.apache.commons.lang3.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestHandler;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.steverhoton.poc.exception.MemberNotFoundException;
import com.steverhoton.poc.exception.StytchApiException;
import com.steverhoton.poc.exception.ValidationException;
import com.steverhoton.poc.model.ApiGatewayRequest;
import com.steverhoton.poc.model.ApiGatewayResponse;
import com.steverhoton.poc.model.ErrorResponse;
import com.steverhoton.poc.model.MemberRequest;
import com.steverhoton.poc.service.StytchMemberService;
import com.steverhoton.poc.stytch.mock.StytchMockModels;

public class MemberLambdaHandler implements RequestHandler<ApiGatewayRequest, ApiGatewayResponse> {

  private static final Logger logger = LoggerFactory.getLogger(MemberLambdaHandler.class);
  private static final ObjectMapper objectMapper = new ObjectMapper();
  private final StytchMemberService memberService;

  public MemberLambdaHandler() {
    String projectId = System.getenv("STYTCH_PROJECT_ID");
    String secret = System.getenv("STYTCH_SECRET");
    String organizationId = System.getenv("STYTCH_ORGANIZATION_ID");

    if (StringUtils.isBlank(projectId)) {
      throw new IllegalStateException("STYTCH_PROJECT_ID environment variable is not set");
    }
    if (StringUtils.isBlank(secret)) {
      throw new IllegalStateException("STYTCH_SECRET environment variable is not set");
    }
    if (StringUtils.isBlank(organizationId)) {
      organizationId = "default-org";
      logger.warn("STYTCH_ORGANIZATION_ID not set, using default: {}", organizationId);
    }

    this.memberService = new StytchMemberService(projectId, secret, organizationId);
    logger.info("MemberLambdaHandler initialized");
  }

  MemberLambdaHandler(StytchMemberService memberService) {
    this.memberService = memberService;
  }

  @Override
  public ApiGatewayResponse handleRequest(ApiGatewayRequest request, Context context) {
    logger.info("Processing request: {} {}", request.getHttpMethod(), request.getPath());

    try {
      String method = request.getHttpMethod();
      String path = request.getPath();

      if (path == null || method == null) {
        return createErrorResponse(400, "Bad Request", "Missing method or path", path);
      }

      if (path.equals("/members") || path.equals("/members/")) {
        return handleMembersEndpoint(request, method);
      } else if (path.matches("/members/[^/]+/?")) {
        return handleMemberByIdEndpoint(request, method, path);
      } else if (path.matches("/members/[^/]+/reactivate/?")) {
        return handleReactivateEndpoint(request, method, path);
      } else {
        return createErrorResponse(404, "Not Found", "Endpoint not found: " + path, path);
      }

    } catch (ValidationException e) {
      logger.error("Validation error: {}", e.getMessage());
      return createErrorResponse(400, "Validation Error", e.getMessage(), request.getPath());
    } catch (MemberNotFoundException e) {
      logger.error("Member not found: {}", e.getMessage());
      return createErrorResponse(404, "Not Found", e.getMessage(), request.getPath());
    } catch (StytchApiException e) {
      logger.error("Stytch API error: {}", e.getMessage());
      int statusCode = e.getStatusCode() == 0 ? 500 : e.getStatusCode();
      return createErrorResponse(statusCode, "API Error", e.getMessage(), request.getPath());
    } catch (Exception e) {
      logger.error("Unexpected error", e);
      return createErrorResponse(
          500,
          "Internal Server Error",
          "An unexpected error occurred: " + e.getMessage(),
          request.getPath());
    }
  }

  private ApiGatewayResponse handleMembersEndpoint(ApiGatewayRequest request, String method)
      throws Exception {
    switch (method) {
      case "POST":
        return handleCreateMember(request);
      case "GET":
        return handleSearchMembers(request);
      default:
        return createErrorResponse(
            405,
            "Method Not Allowed",
            "Method " + method + " not allowed on /members",
            request.getPath());
    }
  }

  private ApiGatewayResponse handleMemberByIdEndpoint(
      ApiGatewayRequest request, String method, String path) throws Exception {
    String memberId = extractMemberId(path);
    if (StringUtils.isBlank(memberId)) {
      return createErrorResponse(400, "Bad Request", "Invalid member ID", path);
    }

    switch (method) {
      case "GET":
        return handleGetMember(memberId);
      case "PUT":
        return handleUpdateMember(request, memberId);
      case "DELETE":
        return handleDeleteMember(memberId);
      default:
        return createErrorResponse(
            405, "Method Not Allowed", "Method " + method + " not allowed on " + path, path);
    }
  }

  private ApiGatewayResponse handleReactivateEndpoint(
      ApiGatewayRequest request, String method, String path) throws Exception {
    if (!"POST".equals(method)) {
      return createErrorResponse(
          405, "Method Not Allowed", "Method " + method + " not allowed on " + path, path);
    }

    String memberId = extractMemberIdFromReactivatePath(path);
    if (StringUtils.isBlank(memberId)) {
      return createErrorResponse(400, "Bad Request", "Invalid member ID", path);
    }

    return handleReactivateMember(memberId);
  }

  private ApiGatewayResponse handleCreateMember(ApiGatewayRequest request) throws Exception {
    String body = request.getBody();
    if (StringUtils.isBlank(body)) {
      throw new ValidationException("Request body is required");
    }

    MemberRequest memberRequest = objectMapper.readValue(body, MemberRequest.class);
    StytchMockModels.Member member = memberService.createMember(memberRequest);
    String responseBody = objectMapper.writeValueAsString(member);
    return ApiGatewayResponse.created(responseBody);
  }

  private ApiGatewayResponse handleGetMember(String memberId) throws Exception {
    StytchMockModels.Member member = memberService.getMember(memberId);
    String responseBody = objectMapper.writeValueAsString(member);
    return ApiGatewayResponse.ok(responseBody);
  }

  private ApiGatewayResponse handleSearchMembers(ApiGatewayRequest request) throws Exception {
    Map<String, String> queryParams = request.getQueryStringParameters();
    StytchMockModels.SearchResponse response = memberService.searchMembers(queryParams);
    String responseBody = objectMapper.writeValueAsString(response);
    return ApiGatewayResponse.ok(responseBody);
  }

  private ApiGatewayResponse handleUpdateMember(ApiGatewayRequest request, String memberId)
      throws Exception {
    String body = request.getBody();
    if (StringUtils.isBlank(body)) {
      throw new ValidationException("Request body is required");
    }

    MemberRequest memberRequest = objectMapper.readValue(body, MemberRequest.class);
    StytchMockModels.Member member = memberService.updateMember(memberId, memberRequest);
    String responseBody = objectMapper.writeValueAsString(member);
    return ApiGatewayResponse.ok(responseBody);
  }

  private ApiGatewayResponse handleDeleteMember(String memberId) {
    memberService.deleteMember(memberId);
    return ApiGatewayResponse.noContent();
  }

  private ApiGatewayResponse handleReactivateMember(String memberId) throws Exception {
    StytchMockModels.Member member = memberService.reactivateMember(memberId);
    String responseBody = objectMapper.writeValueAsString(member);
    return ApiGatewayResponse.ok(responseBody);
  }

  private String extractMemberId(String path) {
    String[] parts = path.split("/");
    if (parts.length >= 3) {
      return parts[2];
    }
    return null;
  }

  private String extractMemberIdFromReactivatePath(String path) {
    String[] parts = path.split("/");
    if (parts.length >= 3) {
      return parts[2];
    }
    return null;
  }

  private ApiGatewayResponse createErrorResponse(
      int statusCode, String error, String message, String path) {
    try {
      ErrorResponse errorResponse = new ErrorResponse(error, message, path);
      String body = objectMapper.writeValueAsString(errorResponse);
      return new ApiGatewayResponse(statusCode, body);
    } catch (Exception e) {
      logger.error("Failed to create error response", e);
      return new ApiGatewayResponse(
          statusCode, "{\"error\":\"" + error + "\",\"message\":\"" + message + "\"}");
    }
  }
}
