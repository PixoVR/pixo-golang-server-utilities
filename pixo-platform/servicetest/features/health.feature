
Feature: Acceptance Tests

  Scenario: Basic health check
    When I send "GET" request to "/health"
    Then the response code should be "200"
    And the response should contain a "ok"
    And the response should contain a "$STATIC_VAL"
    And the response should contain a "$DYNAMIC_VAL"

  Scenario: Custom steps
    And I can say hello
    And I can say goodbye
