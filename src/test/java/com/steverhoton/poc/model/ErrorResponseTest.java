package com.steverhoton.poc.model;

import static org.assertj.core.api.Assertions.assertThat;

import java.time.Instant;

import org.junit.jupiter.api.Test;

class ErrorResponseTest {

  @Test
  void constructor_setsAllFields() {
    ErrorResponse response = new ErrorResponse("NotFound", "Resource not found", "/api/test");

    assertThat(response.getError()).isEqualTo("NotFound");
    assertThat(response.getMessage()).isEqualTo("Resource not found");
    assertThat(response.getPath()).isEqualTo("/api/test");
    assertThat(response.getTimestamp()).isNotNull();

    Instant timestamp = Instant.parse(response.getTimestamp());
    assertThat(timestamp).isBeforeOrEqualTo(Instant.now());
  }

  @Test
  void settersAndGetters_workCorrectly() {
    ErrorResponse response = new ErrorResponse("Error", "Message", "/path");

    response.setError("NewError");
    response.setMessage("NewMessage");
    response.setPath("/newpath");
    response.setTimestamp("2023-01-01T00:00:00Z");

    assertThat(response.getError()).isEqualTo("NewError");
    assertThat(response.getMessage()).isEqualTo("NewMessage");
    assertThat(response.getPath()).isEqualTo("/newpath");
    assertThat(response.getTimestamp()).isEqualTo("2023-01-01T00:00:00Z");
  }
}
