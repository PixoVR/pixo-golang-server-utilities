
Feature: Acceptance Tests

  Scenario: Basic health check
    When I send "GET" request to "/api/health"
    Then the response code should be "200"
    And the response should contain a "ok"

  Scenario: Custom steps
    And I can say hello
    And I can say goodbye
