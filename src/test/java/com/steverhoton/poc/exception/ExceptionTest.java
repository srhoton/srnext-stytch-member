package com.steverhoton.poc.exception;

import static org.assertj.core.api.Assertions.assertThat;

import org.junit.jupiter.api.Test;

class ExceptionTest {

  @Test
  void memberNotFoundException_constructors() {
    MemberNotFoundException ex1 = new MemberNotFoundException("Member not found");
    assertThat(ex1.getMessage()).isEqualTo("Member not found");

    Exception cause = new Exception("cause");
    MemberNotFoundException ex2 = new MemberNotFoundException("Member not found", cause);
    assertThat(ex2.getMessage()).isEqualTo("Member not found");
    assertThat(ex2.getCause()).isEqualTo(cause);
  }

  @Test
  void validationException_constructors() {
    ValidationException ex1 = new ValidationException("Validation failed");
    assertThat(ex1.getMessage()).isEqualTo("Validation failed");

    Exception cause = new Exception("cause");
    ValidationException ex2 = new ValidationException("Validation failed", cause);
    assertThat(ex2.getMessage()).isEqualTo("Validation failed");
    assertThat(ex2.getCause()).isEqualTo(cause);
  }

  @Test
  void stytchApiException_constructors() {
    StytchApiException ex1 = new StytchApiException("API Error", 400);
    assertThat(ex1.getMessage()).isEqualTo("API Error");
    assertThat(ex1.getStatusCode()).isEqualTo(400);

    Exception cause = new Exception("cause");
    StytchApiException ex2 = new StytchApiException("API Error", 500, cause);
    assertThat(ex2.getMessage()).isEqualTo("API Error");
    assertThat(ex2.getStatusCode()).isEqualTo(500);
    assertThat(ex2.getCause()).isEqualTo(cause);
  }
}
