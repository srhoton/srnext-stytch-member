package com.steverhoton.poc.exception;

public class StytchApiException extends RuntimeException {

  private final int statusCode;

  public StytchApiException(String message, int statusCode) {
    super(message);
    this.statusCode = statusCode;
  }

  public StytchApiException(String message, int statusCode, Throwable cause) {
    super(message, cause);
    this.statusCode = statusCode;
  }

  public int getStatusCode() {
    return statusCode;
  }
}
