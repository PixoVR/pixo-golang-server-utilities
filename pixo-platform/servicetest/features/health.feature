
Feature: Acceptance Tests

  Scenario: Basic health check
    When I send "GET" request to "/health"
    Then the response code should be "200"
    And the response should contain a "ok"
    And the response should contain a "$STATIC_VAL"
    And the response should contain a "$DYNAMIC_VAL"
    And the response should contain a "$CUSTOM_VAL"


  Scenario: External health check
    When I send "GET" request to the "allocator" service at "/health"
    Then the response code should be "200"


  Scenario: Not found
    When I send "GET" request to "/nonexistent"
    Then the response code should be "404"
    And the response should not contain a "ok"


  Scenario: Custom steps
    And I can say hello
    And I can say goodbye
