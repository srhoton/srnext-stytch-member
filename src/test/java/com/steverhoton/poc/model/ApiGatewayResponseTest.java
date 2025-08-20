package com.steverhoton.poc.model;

import static org.assertj.core.api.Assertions.assertThat;

import org.junit.jupiter.api.Test;

class ApiGatewayResponseTest {

  @Test
  void constructor_setsDefaultHeaders() {
    ApiGatewayResponse response = new ApiGatewayResponse();

    assertThat(response.getHeaders()).containsEntry("Content-Type", "application/json");
    assertThat(response.getHeaders()).containsEntry("Access-Control-Allow-Origin", "*");
    assertThat(response.getHeaders())
        .containsEntry("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS");
    assertThat(response.getHeaders())
        .containsEntry("Access-Control-Allow-Headers", "Content-Type, Authorization");
  }

  @Test
  void constructor_withStatusAndBody() {
    ApiGatewayResponse response = new ApiGatewayResponse(200, "test body");

    assertThat(response.getStatusCode()).isEqualTo(200);
    assertThat(response.getBody()).isEqualTo("test body");
    assertThat(response.isBase64Encoded()).isFalse();
  }

  @Test
  void ok_returnsCorrectResponse() {
    ApiGatewayResponse response = ApiGatewayResponse.ok("success");

    assertThat(response.getStatusCode()).isEqualTo(200);
    assertThat(response.getBody()).isEqualTo("success");
  }

  @Test
  void created_returnsCorrectResponse() {
    ApiGatewayResponse response = ApiGatewayResponse.created("created");

    assertThat(response.getStatusCode()).isEqualTo(201);
    assertThat(response.getBody()).isEqualTo("created");
  }

  @Test
  void noContent_returnsCorrectResponse() {
    ApiGatewayResponse response = ApiGatewayResponse.noContent();

    assertThat(response.getStatusCode()).isEqualTo(204);
    assertThat(response.getBody()).isNull();
  }

  @Test
  void badRequest_returnsCorrectResponse() {
    ApiGatewayResponse response = ApiGatewayResponse.badRequest("bad request");

    assertThat(response.getStatusCode()).isEqualTo(400);
    assertThat(response.getBody()).isEqualTo("bad request");
  }

  @Test
  void notFound_returnsCorrectResponse() {
    ApiGatewayResponse response = ApiGatewayResponse.notFound("not found");

    assertThat(response.getStatusCode()).isEqualTo(404);
    assertThat(response.getBody()).isEqualTo("not found");
  }

  @Test
  void internalServerError_returnsCorrectResponse() {
    ApiGatewayResponse response = ApiGatewayResponse.internalServerError("error");

    assertThat(response.getStatusCode()).isEqualTo(500);
    assertThat(response.getBody()).isEqualTo("error");
  }

  @Test
  void settersAndGetters_workCorrectly() {
    ApiGatewayResponse response = new ApiGatewayResponse();

    response.setStatusCode(201);
    response.setBody("test");
    response.setBase64Encoded(true);

    assertThat(response.getStatusCode()).isEqualTo(201);
    assertThat(response.getBody()).isEqualTo("test");
    assertThat(response.isBase64Encoded()).isTrue();
  }
}
