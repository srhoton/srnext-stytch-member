package com.steverhoton.poc.model;

import java.util.HashMap;
import java.util.Map;

import com.fasterxml.jackson.annotation.JsonProperty;

public class ApiGatewayResponse {

  @JsonProperty("statusCode")
  private int statusCode;

  @JsonProperty("headers")
  private Map<String, String> headers;

  @JsonProperty("body")
  private String body;

  @JsonProperty("isBase64Encoded")
  private boolean isBase64Encoded;

  public ApiGatewayResponse() {
    this.headers = new HashMap<>();
    this.headers.put("Content-Type", "application/json");
    this.headers.put("Access-Control-Allow-Origin", "*");
    this.headers.put("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS");
    this.headers.put("Access-Control-Allow-Headers", "Content-Type, Authorization");
  }

  public ApiGatewayResponse(int statusCode, String body) {
    this();
    this.statusCode = statusCode;
    this.body = body;
    this.isBase64Encoded = false;
  }

  public static ApiGatewayResponse ok(String body) {
    return new ApiGatewayResponse(200, body);
  }

  public static ApiGatewayResponse created(String body) {
    return new ApiGatewayResponse(201, body);
  }

  public static ApiGatewayResponse noContent() {
    return new ApiGatewayResponse(204, null);
  }

  public static ApiGatewayResponse badRequest(String body) {
    return new ApiGatewayResponse(400, body);
  }

  public static ApiGatewayResponse notFound(String body) {
    return new ApiGatewayResponse(404, body);
  }

  public static ApiGatewayResponse internalServerError(String body) {
    return new ApiGatewayResponse(500, body);
  }

  public int getStatusCode() {
    return statusCode;
  }

  public void setStatusCode(int statusCode) {
    this.statusCode = statusCode;
  }

  public Map<String, String> getHeaders() {
    return headers;
  }

  public void setHeaders(Map<String, String> headers) {
    this.headers = headers;
  }

  public String getBody() {
    return body;
  }

  public void setBody(String body) {
    this.body = body;
  }

  public boolean isBase64Encoded() {
    return isBase64Encoded;
  }

  public void setBase64Encoded(boolean base64Encoded) {
    isBase64Encoded = base64Encoded;
  }
}
