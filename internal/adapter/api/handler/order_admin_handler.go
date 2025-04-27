package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/application/dto"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/mappers"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/ports"

	"github.com/gin-gonic/gin"
)

type OrderAdminHandler struct {
	updateOrderUseCase        usecase.UpdateOrderUseCase
	getPaginatedOrdersUseCase usecase.GetPaginatedOrdersUseCase
}

func NewOrderAdminHandler(
	updateOrderUseCase usecase.UpdateOrderUseCase,
	getPaginatedOrdersUseCase usecase.GetPaginatedOrdersUseCase,
) *OrderAdminHandler {
	return &OrderAdminHandler{
		updateOrderUseCase:        updateOrderUseCase,
		getPaginatedOrdersUseCase: getPaginatedOrdersUseCase,
	}
}

// UpdateOrderStatus endpoit.
// @Summary	    Update order status.
// @Description Update order status.
// @Tags        Order admin
// @Accept      json
// @Produce     json
// @Param       orderID path string true "order id"
// @Param       status path string true "order status"
// @Success     204 "No Content"
// @Failure     400 {object} dto.APIErrorsOutput
// @Failure     500 {object} dto.APIErrorsOutput
// @Router      /api/v1/admin/{orderID}/{status} [patch].
func (oah *OrderAdminHandler) UpdateOrderStatus(c *gin.Context) {
	ctx := c.Request.Context()

	orderID, err := strconv.Atoi(c.Param("orderID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.SimpleAPIErrorsOutput(
			"",
			"orderID",
			"Invalid order ID",
		))
		return
	}

	status := c.Param("status")
	if status != "ready" && status != "delivered" {
		c.JSON(http.StatusBadRequest, dto.SimpleAPIErrorsOutput(
			"",
			"status",
			"Invalid order status",
		))
		return
	}

	if err = oah.updateOrderUseCase.Run(ctx, uint(orderID), status); err != nil { //nolint:gosec // not necessary overflow.
		if errors.Is(err, datastore.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, dto.SimpleAPIErrorsOutput(
				"",
				"orderID",
				"Order not found",
			))
			return
		}

		c.JSON(http.StatusInternalServerError, dto.SimpleAPIErrorsOutput(
			"",
			"",
			"Internal server error",
		))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetPaginatedOrders in handler.
// @Summary	    Get paginated orders.
// @Description Get paginated orders.
// @Tags        Order admin
// @Accept      json
// @Produce     json
// @Param       page query string false "page number"
// @Param       pageSize query string false "page size"
// @Success     200 {object} []dto.ClientOutput
// @Failure     500 {object} dto.APIErrorsOutput
// @Router      /api/v1/admin/orders [get].
func (oah *OrderAdminHandler) GetPaginatedOrders(c *gin.Context) {
	ctx := c.Request.Context()

	// Corrigir capturação de page e pageSize como parâmetros de query
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.SimpleAPIErrorsOutput(
			"",
			"page",
			"Invalid page parameter",
		))
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.SimpleAPIErrorsOutput(
			"",
			"pageSize",
			"Invalid pageSize parameter",
		))
		return
	}

	filter := ports.Filter{
		Limit:  pageSize,
		Offset: page * pageSize,
	}

	orders, err := oah.getPaginatedOrdersUseCase.Run(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.SimpleAPIErrorsOutput(
			"",
			"",
			"Failed to get orders",
		))
		return
	}

	ordersDTO := make([]dto.OrderOutput, len(orders))
	for i, order := range orders {
		ordersDTO[i] = mappers.ToOrderDTO(*order)
	}

	c.JSON(http.StatusOK, ordersDTO)
}
