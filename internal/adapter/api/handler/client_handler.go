package handler

import (
	"errors"
	"net/http"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/datastore"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/application/dto"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase"
	useCaseDTO "github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/dto"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase/mappers"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/shared/validation"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	createClientUseCase usecase.CreateClientUseCase
	getClientUseCase    usecase.GetClientUseCase
}

func NewClientHandler(
	createClientUseCase usecase.CreateClientUseCase,
	getClientUseCase usecase.GetClientUseCase,
) *ClientHandler {
	return &ClientHandler{
		createClientUseCase: createClientUseCase,
		getClientUseCase:    getClientUseCase,
	}
}

// Create client.
// @Summary	    Create client.
// @Description Create client.
// @Tags        Client
// @Accept      json
// @Produce     json
// @Param       client body useCaseDTO.ClientInputCreate true "request body"
// @Success     201 {object} dto.ClientOutput
// @Failure     400 {object} dto.APIErrorsOutput
// @Failure     409 {object} dto.APIErrorsOutput
// @Failure     500 {object} dto.APIErrorsOutput
// @Router      /api/v1/clients/ [post].
func (ch *ClientHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var input useCaseDTO.ClientInputCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.SimpleAPIErrorsOutput(
			"",
			"body",
			"Invalid request body",
		))
		return
	}

	if err := useCaseDTO.ValidateClientCreate(input); err != nil {
		errors := validation.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, dto.ErrorsFromValidationErrors(errors))
		return
	}

	client, err := ch.createClientUseCase.Run(ctx, input)
	if err != nil {
		if errors.Is(err, datastore.ErrExistingRecord) {
			c.JSON(http.StatusConflict, dto.SimpleAPIErrorsOutput(
				"",
				"cpf",
				"cpf already exists",
			))
			return
		}
		c.JSON(http.StatusInternalServerError, dto.SimpleAPIErrorsOutput("", "", "failed to create client"))
		return
	}

	c.JSON(http.StatusCreated, mappers.ToClientDTO(*client))
}

// GetClient by CPF.
// @Summary	    Get client by CPF.
// @Description Get client by CPF.
// @Tags        Client
// @Accept      json
// @Produce     json
// @Param       cpf path string true "client cpf"
// @Success     200 {object} dto.ClientOutput
// @Failure     404 "No Content"
// @Failure     500 {object} dto.APIErrorsOutput
// @Router      /api/v1/clients/{cpf} [get].
func (ch *ClientHandler) GetClient(c *gin.Context) {
	ctx := c.Request.Context()
	cpf := c.Param("cpf")

	client, err := ch.getClientUseCase.Run(ctx, cpf)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, dto.SimpleAPIErrorsOutput(
				"",
				"",
				"Failed to get client",
			))
			return
		}
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, mappers.ToClientDTO(*client))
}
