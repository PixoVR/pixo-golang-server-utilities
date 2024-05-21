
@health
Feature: Health Check

  Scenario:
    When I send "GET" request to "/api/health"
    Then the response code should be "200"
    And the response should contain a "ok"
