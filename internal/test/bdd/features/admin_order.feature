Feature: Admin order management

  Scenario: Update order status to ready successfully
    Given an existing order with status "pending"
    When I PATCH "/api/v1/admin/orders/{orderID}/ready"
    Then the response status should be 204
    And the order status in the database should be "ready"

  Scenario: Update order status with invalid ID
    When I PATCH "/api/v1/admin/orders/invalid-id/ready"
    Then the response status should be 400
    And the response should contain "Invalid order ID"

  Scenario: Update order status with invalid status
    When I PATCH "/api/v1/admin/orders/1/invalid-status"
    Then the response status should be 400
    And the response should contain "Invalid order status"

  Scenario: Update order status with non-existing ID
    When I PATCH "/api/v1/admin/orders/99999/ready"
    Then the response status should be 404
    And the response should contain "Order not found"
