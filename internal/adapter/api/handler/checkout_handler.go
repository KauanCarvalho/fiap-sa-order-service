package handler

import (
	"errors"
	"net/http"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/clients/payment"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/clients/product"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/application/dto"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase"
	useCaseDTO "github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/dto"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/mappers"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/shared/validation"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type CheckoutHandler struct {
	createOrderUseCase usecase.CreateOrderUseCase
}

func NewCheckoutHandler(createOrderUseCase usecase.CreateOrderUseCase) *CheckoutHandler {
	return &CheckoutHandler{
		createOrderUseCase: createOrderUseCase,
	}
}

// Create order.
// @Summary	    Create order.
// @Description Create order.
// @Tags        Checkout
// @Accept      json
// @Produce     json
// @Param       order body useCaseDTO.OrderInputCreate true "request body"
// @Success     201 {object} dto.OrderOutput
// @Failure     400 {object} dto.APIErrorsOutput
// @Failure     409 {object} dto.APIErrorsOutput
// @Failure     500 {object} dto.APIErrorsOutput
// @Router      /api/v1/orders/ [post].
func (ch *CheckoutHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var input useCaseDTO.OrderInputCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.SimpleAPIErrorsOutput(
			"",
			"body",
			"Invalid request body",
		))
		return
	}

	if err := useCaseDTO.ValidateOrderCreate(input); err != nil {
		errors := validation.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, dto.ErrorsFromValidationErrors(errors))
		return
	}

	order, err := ch.createOrderUseCase.Run(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusBadRequest, dto.SimpleAPIErrorsOutput(
				"",
				"client_id",
				"client not found",
			))
		case errors.Is(err, product.ErrSKUNotFound):
			c.JSON(http.StatusBadRequest, dto.SimpleAPIErrorsOutput(
				"",
				"items",
				"product not found",
			))
		case errors.Is(err, payment.ErrProblemToAuthorizePayment):
			c.JSON(http.StatusInternalServerError, dto.SimpleAPIErrorsOutput(
				"",
				"payment",
				"problem to authorize payment",
			))
		default:
			c.JSON(http.StatusInternalServerError, dto.SimpleAPIErrorsOutput("", "", "failed to create order"))
		}
		return
	}

	c.JSON(http.StatusCreated, mappers.ToOrderDTO(*order))
}
