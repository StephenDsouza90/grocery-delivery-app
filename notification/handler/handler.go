package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	k "github.com/StephenDsouza90/grocery-delivery-app/internal/kafka"
	r "github.com/StephenDsouza90/grocery-delivery-app/internal/repository"
	u "github.com/StephenDsouza90/grocery-delivery-app/internal/utils"
)

// Handler provides HTTP handlers for interacting with orders.
type Handler struct {
	repo     r.Repository
	producer *k.Producer
	consumer *k.Consumer
}

// NewHandler creates a new Handler.
func NewHandler(repo r.Repository, consumer *k.Consumer) *Handler {
	return &Handler{
		repo:     repo,
		consumer: consumer,
	}
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
		var baseEvent map[string]interface{}
		if err := json.Unmarshal(msg.Value, &baseEvent); err != nil {
			log.Printf("Failed to unmarshal base event: %v", err)
			continue
		}

		if eventType, ok := baseEvent["type"].(string); ok {
			switch eventType {
			case "order":
				var order r.Order
				if err := json.Unmarshal(msg.Value, &order); err != nil {
					log.Printf("Failed to unmarshal order event: %v", err)
					continue
				}
				if err := h.OrderNotification(context.Background(), order); err != nil {
					log.Printf("Failed to notify user about order: %v", err)
				}
			case "payment":
				var payment r.Payment
				if err := json.Unmarshal(msg.Value, &payment); err != nil {
					log.Printf("Failed to unmarshal payment event: %v", err)
					continue
				}
				if err := h.PaymentNotification(context.Background(), payment); err != nil {
					log.Printf("Failed to notify user about payment: %v", err)
				}
			case "delivery":
				var delivery r.Delivery
				if err := json.Unmarshal(msg.Value, &delivery); err != nil {
					log.Printf("Failed to unmarshal delivery event: %v", err)
					continue
				}
				if err := h.DeliveryNotification(context.Background(), delivery); err != nil {
					log.Printf("Failed to notify user about delivery: %v", err)
				}
			default:
				log.Printf("Unknown event type: %s", eventType)
			}
		} else {
			log.Printf("Event type not found in message")
		}

		session.MarkMessage(msg, "")
	}
	return nil
}

// OrderNotification sends a notification to the user about the order.
func (h *Handler) OrderNotification(c context.Context, o r.Order) error {
	if o.OrderID != "" {
		notification := r.Notification{
			NotificationID: u.GenerateUUID(),
			OrderID:        o.OrderID,
			Message:        fmt.Sprintf("Your order with ID %s has been %s", o.OrderID, o.OrderStatus),
		}

		if err := h.repo.AddNotification(notification); err != nil {
			return err
		}

		mockNotificationSentToUser(notification.Message)
	} else {
		log.Printf("Order ID is empty")
	}

	return nil
}

// PaymentNotification sends a notification to the user about the payment.
func (h *Handler) PaymentNotification(c context.Context, p r.Payment) error {
	if p.PaymentID != "" {
		notification := r.Notification{
			NotificationID: u.GenerateUUID(),
			OrderID:        p.OrderID,
			Message:        fmt.Sprintf("Your payment for order ID %s has been %s", p.OrderID, p.PaymentStatus),
		}

		if err := h.repo.AddNotification(notification); err != nil {
			return err
		}

		mockNotificationSentToUser(notification.Message)

	} else {
		log.Printf("Payment ID is empty")
	}

	return nil
}

// DeliveryNotification sends a notification to the user about the delivery.
func (h *Handler) DeliveryNotification(c context.Context, d r.Delivery) error {
	if d.DeliveryID != "" {
		notification := r.Notification{
			NotificationID: u.GenerateUUID(),
			OrderID:        d.OrderID,
			Message:        fmt.Sprintf("Your order with ID %s has been %s", d.OrderID, d.DeliveryStatus),
		}

		if err := h.repo.AddNotification(notification); err != nil {
			return err
		}

		mockNotificationSentToUser(notification.Message)
	} else {
		log.Printf("Delivery ID is empty")
	}

	return nil
}

func mockNotificationSentToUser(message string) {
	fmt.Println("Notification sent to user: ", message)
	time.Sleep(2 * time.Second)
}
