package delivery

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/StephenDsouza90/grocery-delivery-app/internal/kafka"
	"github.com/StephenDsouza90/grocery-delivery-app/internal/repository"
)

// Handler provides HTTP handlers for interacting with orders.
type Handler struct {
	repo     repository.Repository
	producer *kafka.Producer
	consumer *kafka.Consumer
}

type PaymentCreatedMessage struct {
	OrderID       string `json:"order_id"`
	PaymentStatus string `json:"payment_status"`
}

type OrderCreatedMessage struct {
	OrderID      string `json:"order_id"`
	DeliveryDate string `json:"delivery_date"`
	DeliveryTime string `json:"delivery_time"`
}

// NewHandler creates a new Handler.
func NewHandler(repo repository.Repository, producer *kafka.Producer, consumer *kafka.Consumer) *Handler {
	handler := &Handler{
		repo:     repo,
		producer: producer,
		consumer: consumer,
	}
	return handler
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *Handler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *Handler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim consumes messages from the Kafka topic.
func (h *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var payment PaymentCreatedMessage
		if err := json.Unmarshal(msg.Value, &payment); err != nil {
			log.Printf("Failed to unmarshal PaymentCreatedMessage event: %v", err)
			continue
		}

		var order OrderCreatedMessage
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Failed to unmarshal OrderCreatedMessage event: %v", err)
			continue
		}

		// Assign delivery
		if err := h.AssignDelivery(context.Background(), payment, order); err != nil {
			log.Printf("Payment processing failed: %v", err)
		}

		session.MarkMessage(msg, "")
	}
	return nil
}

// AssignDelivery assigns a delivery to an order.
func (h *Handler) AssignDelivery(ctx context.Context, payment PaymentCreatedMessage, order OrderCreatedMessage) error {
	// Check if payment was successful
	if payment.PaymentStatus != "success" {
		log.Printf("Payment was not successful for order %s", order.OrderID)
	}

	// Get delivery date and time and assign it for delivery
	deliveryDate := order.DeliveryDate
	deliveryTime := order.DeliveryTime

	// Mock delivery assignment
	deliveryStatus := mockDelivery(deliveryDate, deliveryTime)

	orderID, err := strconv.Atoi(order.OrderID)
	if err != nil {
		log.Printf("Failed to convert order ID to int: %v", err)
	}

	// Save delivery to database
	delivery := repository.Delivery{
		DeliveryID:       generateDeliveryId(),
		OrderID:          orderID,
		DeliveryPersonID: 1,
		DeliveryStatus:   deliveryStatus,
	}

	// Save the order
	if err := h.repo.AddDelivery(delivery); err != nil {
		log.Printf("Failed to save delivery: %v", err)
		return err
	}

	// Update status of order
	var orderStatus string
	if deliveryStatus == "assigned" {
		orderStatus = "assigned"
	} else {
		orderStatus = "pending"
	}
	if err := h.repo.UpdateStatus(orderID, orderStatus); err != nil {
		log.Printf("Failed to update order status: %v", err)
	}

	// Prepare the delivery data
	deliveryData := map[string]string{
		"order_id":        order.OrderID,
		"delivery_id":     strconv.Itoa(delivery.DeliveryID),
		"delivery_status": deliveryStatus,
		"delivery_date":   deliveryDate,
		"delivery_time":   deliveryTime,
	}

	// Send delivery message
	h.producer.SendDeliveryMessage(deliveryData)

	return nil
}

// mockDelivery mocks the delivery assignment.
func mockDelivery(deliveryDate string, deliveryTime string) string {
	if deliveryDate == "" && deliveryTime == "" {
		return "pending"
	}
	return "assigned"
}

func generateDeliveryId() int {
	return int(time.Now().Unix())
}
