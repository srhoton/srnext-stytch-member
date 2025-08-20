package com.steverhoton.poc.model;

import java.time.Instant;

import com.fasterxml.jackson.annotation.JsonProperty;

public class ErrorResponse {

  @JsonProperty("error")
  private String error;

  @JsonProperty("message")
  private String message;

  @JsonProperty("timestamp")
  private String timestamp;

  @JsonProperty("path")
  private String path;

  public ErrorResponse(String error, String message, String path) {
    this.error = error;
    this.message = message;
    this.timestamp = Instant.now().toString();
    this.path = path;
  }

  public String getError() {
    return error;
  }

  public void setError(String error) {
    this.error = error;
  }

  public String getMessage() {
    return message;
  }

  public void setMessage(String message) {
    this.message = message;
  }

  public String getTimestamp() {
    return timestamp;
  }

  public void setTimestamp(String timestamp) {
    this.timestamp = timestamp;
  }

  public String getPath() {
    return path;
  }

  public void setPath(String path) {
    this.path = path;
  }
}
