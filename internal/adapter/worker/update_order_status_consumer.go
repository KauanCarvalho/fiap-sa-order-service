package worker

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/application/dto"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/usecase"
	"go.uber.org/zap"
)

type UpdateOrderStatusConnsumer struct {
	updateOrderUseCase usecase.UpdateOrderUseCase
}

func NewUpdateOrderStatusConnsumer(updateOrderUseCase usecase.UpdateOrderUseCase) *UpdateOrderStatusConnsumer {
	return &UpdateOrderStatusConnsumer{
		updateOrderUseCase: updateOrderUseCase,
	}
}

func (consumer UpdateOrderStatusConnsumer) Process(_ context.Context, processingMessage ProcessingMessage) error {
	parsedBody := dto.UpdateOrderStatusInput{}

	envlope := SNSEnvelope{}
	err := json.Unmarshal([]byte(*processingMessage.Message.Body), &envlope)
	if err != nil {
		zap.L().Error(
			"error parsing sns envelope",
			zap.String("queueName", processingMessage.QueueName),
			zap.Error(err),
		)
		return err
	}

	err = json.Unmarshal([]byte(envlope.Message), &parsedBody)
	if err != nil {
		zap.L().Error(
			"error parsing message body",
			zap.String("queueName", processingMessage.QueueName),
			zap.Error(err),
		)
		return err
	}

	ctx := context.Background()

	uint64ID, err := strconv.ParseUint(parsedBody.ExternalRef, 10, 64)
	if err != nil {
		zap.L().Error(
			"error converting external reference to uint",
			zap.String("queueName", processingMessage.QueueName),
			zap.Error(err),
		)
		return err
	}

	err = consumer.updateOrderUseCase.Run(ctx, uint(uint64ID), parsedBody.Status)
	if err != nil {
		zap.L().Error(
			"error updating order status",
			zap.String("queueName", processingMessage.QueueName),
			zap.Error(err),
		)
		return err
	}

	zap.L().Info(
		"order status updated successfully",
		zap.String("queueName", processingMessage.QueueName),
		zap.Uint64("externalReference", uint64ID),
	)
	return nil
}
