package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labbs/nexo/application/webhook/dto"
	"github.com/labbs/nexo/domain"
	"github.com/labbs/nexo/infrastructure/config"
	"github.com/rs/zerolog"
)

type WebhookApp struct {
	Config              config.Config
	Logger              zerolog.Logger
	WebhookPers         domain.WebhookPers
	WebhookDeliveryPers domain.WebhookDeliveryPers
}

func NewWebhookApp(config config.Config, logger zerolog.Logger, webhookPers domain.WebhookPers, webhookDeliveryPers domain.WebhookDeliveryPers) *WebhookApp {
	return &WebhookApp{
		Config:              config,
		Logger:              logger,
		WebhookPers:         webhookPers,
		WebhookDeliveryPers: webhookDeliveryPers,
	}
}

// generateSecret generates a secure random secret for webhook signature
func generateSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "whsec_" + hex.EncodeToString(bytes), nil
}

// signPayload creates an HMAC-SHA256 signature of the payload
func signPayload(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

func (app *WebhookApp) CreateWebhook(input dto.CreateWebhookInput) (*dto.CreateWebhookOutput, error) {
	secret, err := generateSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret: %w", err)
	}

	webhook := &domain.Webhook{
		Id:      uuid.New().String(),
		UserId:  input.UserId,
		SpaceId: input.SpaceId,
		Name:    input.Name,
		Url:     input.Url,
		Secret:  secret,
		Events: domain.JSONB{
			"events": input.Events,
		},
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := app.WebhookPers.Create(webhook); err != nil {
		return nil, fmt.Errorf("failed to create webhook: %w", err)
	}

	return &dto.CreateWebhookOutput{
		Id:     webhook.Id,
		Name:   webhook.Name,
		Url:    webhook.Url,
		Secret: secret,
		Events: input.Events,
		Active: webhook.Active,
	}, nil
}

func (app *WebhookApp) ListWebhooks(input dto.ListWebhooksInput) (*dto.ListWebhooksOutput, error) {
	webhooks, err := app.WebhookPers.GetByUserId(input.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to list webhooks: %w", err)
	}

	output := &dto.ListWebhooksOutput{
		Webhooks: make([]dto.WebhookItem, len(webhooks)),
	}

	for i, w := range webhooks {
		var events []string
		if w.Events != nil {
			if e, ok := w.Events["events"].([]interface{}); ok {
				for _, ev := range e {
					if str, ok := ev.(string); ok {
						events = append(events, str)
					}
				}
			}
		}

		item := dto.WebhookItem{
			Id:           w.Id,
			Name:         w.Name,
			Url:          w.Url,
			SpaceId:      w.SpaceId,
			Events:       events,
			Active:       w.Active,
			LastError:    w.LastError,
			LastErrorAt:  w.LastErrorAt,
			SuccessCount: w.SuccessCount,
			FailureCount: w.FailureCount,
			CreatedAt:    w.CreatedAt,
		}

		if w.Space != nil {
			item.SpaceName = &w.Space.Name
		}

		output.Webhooks[i] = item
	}

	return output, nil
}

func (app *WebhookApp) GetWebhook(input dto.GetWebhookInput) (*dto.GetWebhookOutput, error) {
	webhook, err := app.WebhookPers.GetById(input.WebhookId)
	if err != nil {
		return nil, fmt.Errorf("webhook not found: %w", err)
	}

	if webhook.UserId != input.UserId {
		return nil, fmt.Errorf("access denied")
	}

	var events []string
	if webhook.Events != nil {
		if e, ok := webhook.Events["events"].([]interface{}); ok {
			for _, ev := range e {
				if str, ok := ev.(string); ok {
					events = append(events, str)
				}
			}
		}
	}

	output := &dto.GetWebhookOutput{
		Id:           webhook.Id,
		Name:         webhook.Name,
		Url:          webhook.Url,
		Secret:       webhook.Secret,
		SpaceId:      webhook.SpaceId,
		Events:       events,
		Active:       webhook.Active,
		LastError:    webhook.LastError,
		LastErrorAt:  webhook.LastErrorAt,
		SuccessCount: webhook.SuccessCount,
		FailureCount: webhook.FailureCount,
		CreatedAt:    webhook.CreatedAt,
		UpdatedAt:    webhook.UpdatedAt,
	}

	if webhook.Space != nil {
		output.SpaceName = &webhook.Space.Name
	}

	return output, nil
}

func (app *WebhookApp) UpdateWebhook(input dto.UpdateWebhookInput) error {
	webhook, err := app.WebhookPers.GetById(input.WebhookId)
	if err != nil {
		return fmt.Errorf("webhook not found: %w", err)
	}

	if webhook.UserId != input.UserId {
		return fmt.Errorf("access denied")
	}

	if input.Name != nil {
		webhook.Name = *input.Name
	}

	if input.Url != nil {
		webhook.Url = *input.Url
	}

	if input.Events != nil {
		webhook.Events = domain.JSONB{
			"events": *input.Events,
		}
	}

	if input.Active != nil {
		webhook.Active = *input.Active
	}

	webhook.UpdatedAt = time.Now()

	if err := app.WebhookPers.Update(webhook); err != nil {
		return fmt.Errorf("failed to update webhook: %w", err)
	}

	return nil
}

func (app *WebhookApp) DeleteWebhook(input dto.DeleteWebhookInput) error {
	webhook, err := app.WebhookPers.GetById(input.WebhookId)
	if err != nil {
		return fmt.Errorf("webhook not found: %w", err)
	}

	if webhook.UserId != input.UserId {
		return fmt.Errorf("access denied")
	}

	if err := app.WebhookPers.Delete(input.WebhookId); err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}

	return nil
}

func (app *WebhookApp) GetDeliveries(input dto.GetDeliveriesInput) (*dto.GetDeliveriesOutput, error) {
	// Verify ownership
	webhook, err := app.WebhookPers.GetById(input.WebhookId)
	if err != nil {
		return nil, fmt.Errorf("webhook not found: %w", err)
	}

	if webhook.UserId != input.UserId {
		return nil, fmt.Errorf("access denied")
	}

	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}

	deliveries, err := app.WebhookDeliveryPers.GetByWebhookId(input.WebhookId, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get deliveries: %w", err)
	}

	output := &dto.GetDeliveriesOutput{
		Deliveries: make([]dto.DeliveryItem, len(deliveries)),
	}

	for i, d := range deliveries {
		output.Deliveries[i] = dto.DeliveryItem{
			Id:         d.Id,
			Event:      d.Event,
			StatusCode: d.StatusCode,
			Success:    d.Success,
			Duration:   d.Duration,
			CreatedAt:  d.CreatedAt,
		}
	}

	return output, nil
}

// TriggerWebhooks triggers all matching webhooks for an event
func (app *WebhookApp) TriggerWebhooks(input dto.TriggerWebhookInput) {
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

func (app *WebhookApp) deliverWebhook(webhook domain.Webhook, event string, payload map[string]interface{}) {
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

func (app *WebhookApp) recordDelivery(webhookId, event string, payload map[string]interface{}, statusCode int, response string, success bool, duration int) {
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
