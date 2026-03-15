package dto

type GetDeliveriesInput struct {
	UserId    string
	WebhookId string
	Limit     int
}

type GetDeliveriesOutput struct {
	Deliveries []DeliveryItem
}
