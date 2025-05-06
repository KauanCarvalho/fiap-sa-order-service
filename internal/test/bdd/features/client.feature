Feature: Client management

  Scenario: Successfully create a client
    Given I have a client with name "John Doe" and cpf "12345678909"
    When I send a POST request to "/api/v1/clients"
    Then the response code should be 201
    And the response should contain "John Doe" and "12345678909"

  Scenario: Fail to create client with existing CPF
    Given a client with cpf "98765432100" already exists
    When I send a POST request to "/api/v1/clients" with name "Jane Doe" and cpf "98765432100"
    Then the response code should be 409
    And the response should contain "cpf already exists"

  Scenario: Successfully get a client by cpf
    Given a client with name "Alice Smith" and cpf "11122233344" exists
    When I send a GET request to "/api/v1/clients/11122233344"
    Then the response code should be 200
    And the response should contain "Alice Smith" and "11122233344"

  Scenario: Fail to get client with non-existent cpf
    When I send a GET request to "/api/v1/clients/00000000000"
    Then the response code should be 404
