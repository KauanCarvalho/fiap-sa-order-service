Feature: Checkout

  Scenario: Successful checkout
    Given a client with ID "130" exists
    When I send a POST request to "/api/v1/checkout" with body:
      """
      {
        "client_id": 130,
        "items": [
          {"sku": "TESTSKU", "quantity": 2}
        ]
      }
      """
    Then the response status should be 201
    And the response should contain "payment_method" with value "pix"

  Scenario: Invalid request body
    When I send a POST request to "/api/v1/checkout" with body:
      """
      {invalid-json
      """
    Then the response status should be 400

  Scenario: Empty items list
    Given a client with ID "130" exists
    When I send a POST request to "/api/v1/checkout" with body:
      """
      {
        "client_id": 130,
        "items": []
      }
      """
    Then the response status should be 400

  Scenario: Client not found
    When I send a POST request to "/api/v1/checkout" with body:
      """
      {
        "client_id": 99999,
        "items": [{"sku": "TESTSKU", "quantity": 1}]
      }
      """
    Then the response status should be 400

  Scenario: Product not found
    Given a client with ID "130" exists
    When I send a POST request to "/api/v1/checkout" with body:
      """
      {
        "client_id": 130,
        "items": [{"sku": "INVALIDSKU", "quantity": 1}]
      }
      """
    Then the response status should be 400

  Scenario: Product service returns 500
    Given a client with ID "130" exists
    And the product service is returning 500
    When I send a POST request to "/api/v1/checkout" with body:
      """
      {
        "client_id": 130,
        "items": [{"sku": "FORCE500", "quantity": 1}]
      }
      """
    Then the response status should be 500

  Scenario: Payment service returns 500
    Given a client with ID "130" exists
    And the payment service is returning 500
    When I send a POST request to "/api/v1/checkout" with body:
      """
      {
        "client_id": 130,
        "items": [{"sku": "TESTSKU", "quantity": 1}]
      }
      """
    Then the response status should be 500
