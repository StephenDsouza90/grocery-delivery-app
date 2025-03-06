package kafka

const (
	Brokers = "KAFKA_BROKERS"

	// GroupIDs
	OrderGroupID    = "order-group"
	PaymentGroupID  = "payment-group"
	DeliveryGroupID = "delivery-group"

	// Topics
	OrderCreatedTopic   = "OrderCreated"
	PaymentStatusTopic  = "PaymentStatus"
	DeliveryStatusTopic = "DeliveryStatus"
)
