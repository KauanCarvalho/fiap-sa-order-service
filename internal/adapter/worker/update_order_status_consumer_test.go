package worker_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/KauanCarvalho/fiap-sa-order-service/internal/adapter/worker"
	"github.com/KauanCarvalho/fiap-sa-order-service/internal/core/domain/entities"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/require"
)

func TestUpdateOrderStatusConnsumer_Process(t *testing.T) {
	t.Run("should update order status successfully", func(t *testing.T) {
		prepareTestDatabase()

		order := &entities.Order{}
		err := sqlDB.First(&order, 1).Error
		require.NoError(t, err)

		orderPayload := map[string]string{
			"external_reference": strconv.Itoa(int(order.ID)),
			"status":             "confirmed",
		}
		payloadBytes, _ := json.Marshal(orderPayload)

		snsEnvelope := map[string]string{
			"Type":    "Notification",
			"Message": string(payloadBytes),
		}
		envelopeBytes, _ := json.Marshal(snsEnvelope)

		body := string(envelopeBytes)
		msg := &sqs.Message{Body: &body}

		processingMessage := worker.ProcessingMessage{
			Message:   msg,
			QueueName: "test-queue",
		}

		err = uc.Process(ctx, processingMessage)
		require.NoError(t, err)

		reloadedOrder := &entities.Order{}
		err = sqlDB.First(&reloadedOrder, order.ID).Error
		require.NoError(t, err)
		require.Equal(t, "confirmed", reloadedOrder.Status)
	})

	t.Run("should fail when SNS envelope is invalid JSON", func(t *testing.T) {
		body := "{ invalid json }"
		msg := &sqs.Message{Body: &body}

		processingMessage := worker.ProcessingMessage{
			Message:   msg,
			QueueName: "test-queue",
		}

		err := uc.Process(ctx, processingMessage)
		require.Error(t, err)
	})

	t.Run("should fail when Message inside SNS is invalid JSON", func(t *testing.T) {
		snsEnvelope := map[string]string{
			"Type":    "Notification",
			"Message": "{ not valid json",
		}
		envelopeBytes, _ := json.Marshal(snsEnvelope)

		body := string(envelopeBytes)
		msg := &sqs.Message{Body: &body}

		processingMessage := worker.ProcessingMessage{
			Message:   msg,
			QueueName: "test-queue",
		}

		err := uc.Process(ctx, processingMessage)
		require.Error(t, err)
	})

	t.Run("should fail when external_ref is not a number", func(t *testing.T) {
		orderPayload := map[string]string{
			"external_ref": "abc",
			"status":       "cancelled",
		}
		payloadBytes, _ := json.Marshal(orderPayload)

		snsEnvelope := map[string]string{
			"Message": string(payloadBytes),
		}
		envelopeBytes, _ := json.Marshal(snsEnvelope)

		body := string(envelopeBytes)
		msg := &sqs.Message{Body: &body}

		processingMessage := worker.ProcessingMessage{
			Message:   msg,
			QueueName: "test-queue",
		}

		err := uc.Process(ctx, processingMessage)
		require.Error(t, err)
	})

	t.Run("should fail when order does not exist", func(t *testing.T) {
		orderPayload := map[string]string{
			"external_ref": "99999",
			"status":       "shipped",
		}
		payloadBytes, _ := json.Marshal(orderPayload)

		snsEnvelope := map[string]string{
			"Message": string(payloadBytes),
		}
		envelopeBytes, _ := json.Marshal(snsEnvelope)

		body := string(envelopeBytes)
		msg := &sqs.Message{Body: &body}

		processingMessage := worker.ProcessingMessage{
			Message:   msg,
			QueueName: "test-queue",
		}

		err := uc.Process(ctx, processingMessage)
		require.Error(t, err)
	})
}
