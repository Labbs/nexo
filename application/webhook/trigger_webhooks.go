package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labbs/nexo/application/webhook/dto"
	"github.com/labbs/nexo/domain"
)

// TriggerWebhooks triggers all matching webhooks for an event
func (app *WebhookApplication) TriggerWebhooks(input dto.TriggerWebhookInput) {
	logger := app.Logger.With().Str("component", "webhook.trigger").Str("event", input.Event).Logger()

	webhooks, err := app.WebhookPers.GetActiveByEvent(domain.WebhookEvent(input.Event), input.SpaceId)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get webhooks for event")
		return
	}

	for _, webhook := range webhooks {
		go app.deliverWebhook(webhook, input.Event, input.Payload)
	}
}

func (app *WebhookApplication) deliverWebhook(webhook domain.Webhook, event string, payload map[string]interface{}) {
	logger := app.Logger.With().
		Str("component", "webhook.deliver").
		Str("webhook_id", webhook.Id).
		Str("event", event).
		Logger()

	// Build the webhook payload
	webhookPayload := map[string]interface{}{
		"id":        uuid.New().String(),
		"event":     event,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"data":      payload,
	}

	payloadBytes, err := json.Marshal(webhookPayload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to marshal payload")
		return
	}

	// Create signature
	signature := signPayload(payloadBytes, webhook.Secret)

	// Make the HTTP request
	start := time.Now()
	req, err := http.NewRequest("POST", webhook.Url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		logger.Error().Err(err).Msg("failed to create request")
		app.recordDelivery(webhook.Id, event, webhookPayload, 0, "", false, 0)
		_ = app.WebhookPers.RecordFailure(webhook.Id, err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Signature", signature)
	req.Header.Set("X-Webhook-Event", event)
	req.Header.Set("X-Webhook-Delivery", webhookPayload["id"].(string))

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	duration := int(time.Since(start).Milliseconds())

	if err != nil {
		logger.Error().Err(err).Msg("failed to deliver webhook")
		app.recordDelivery(webhook.Id, event, webhookPayload, 0, err.Error(), false, duration)
		_ = app.WebhookPers.RecordFailure(webhook.Id, err.Error())
		return
	}
	defer resp.Body.Close()

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	responseBody := ""
	if !success {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		responseBody = buf.String()
		if len(responseBody) > 500 {
			responseBody = responseBody[:500]
		}
	}

	app.recordDelivery(webhook.Id, event, webhookPayload, resp.StatusCode, responseBody, success, duration)

	if success {
		_ = app.WebhookPers.IncrementSuccess(webhook.Id)
		logger.Debug().Int("status_code", resp.StatusCode).Msg("webhook delivered successfully")
	} else {
		_ = app.WebhookPers.RecordFailure(webhook.Id, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, responseBody))
		logger.Warn().Int("status_code", resp.StatusCode).Msg("webhook delivery failed")
	}
}

func (app *WebhookApplication) recordDelivery(webhookId, event string, payload map[string]interface{}, statusCode int, response string, success bool, duration int) {
	delivery := &domain.WebhookDelivery{
		Id:         uuid.New().String(),
		WebhookId:  webhookId,
		Event:      event,
		Payload:    domain.JSONB(payload),
		StatusCode: statusCode,
		Response:   response,
		Success:    success,
		Duration:   duration,
		CreatedAt:  time.Now(),
	}

	if err := app.WebhookDeliveryPers.Create(delivery); err != nil {
		app.Logger.Error().Err(err).Msg("failed to record webhook delivery")
	}
}
