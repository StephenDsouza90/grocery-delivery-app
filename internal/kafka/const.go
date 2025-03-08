package kafka

const (
	Brokers = "KAFKA_BROKERS"

	// GroupIDs
	OrderGroupID        = "order-group"
	PaymentGroupID      = "payment-group"
	DeliveryGroupID     = "delivery-group"
	NotificationGroupID = "notification-group"

	// Topics
	OrderCreatedTopic   = "OrderCreated"
	PaymentStatusTopic  = "PaymentStatus"
	DeliveryStatusTopic = "DeliveryStatus"
	NotificationTopic   = "Notification"
)
