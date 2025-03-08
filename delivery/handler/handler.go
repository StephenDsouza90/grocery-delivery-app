package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/IBM/sarama"
	k "github.com/StephenDsouza90/grocery-delivery-app/internal/kafka"
	r "github.com/StephenDsouza90/grocery-delivery-app/internal/repository"
	"github.com/gin-gonic/gin"
)

// Handler provides HTTP handlers for interacting with orders.
type Handler struct {
	repo     r.Repository
	producer *k.Producer
	consumer *k.Consumer
}

// NewHandler creates a new Handler.
func NewHandler(repo r.Repository, producer *k.Producer, consumer *k.Consumer) *Handler {
	return &Handler{
		repo:     repo,
		producer: producer,
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
		var payment r.Payment
		if err := json.Unmarshal(msg.Value, &payment); err != nil {
			log.Printf("Failed to unmarshal payment event: %v", err)
			continue
		}

		var order r.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Failed to unmarshal order event: %v", err)
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
func (h *Handler) AssignDelivery(c context.Context, p r.Payment, o r.Order) error {
	if p.PaymentStatus != "success" {
		log.Printf("Payment was not successful for order %d", o.OrderID)
		// TODO : Think what to do if payment is not successful
	} else {
		deliveryDate := o.DeliveryDate
		deliveryTime := o.DeliveryTime

		// Mock delivery assignment
		deliveryStatus := mockDelivery(deliveryDate, deliveryTime)

		// Create and save delivery to database
		newDelivery := deliveryObject(o, deliveryStatus)
		deliveryID, err := h.repo.AddDelivery(newDelivery)
		if err != nil {
			log.Printf("Failed to add delivery for order %d: %v", o.OrderID, err)
		}

		// Send delivery message
		newDelivery.DeliveryID = deliveryID
		h.producer.SendDeliveryMessage(newDelivery)
	}

	return nil
}

// mockDelivery mocks the delivery assignment.
func mockDelivery(deliveryDate string, deliveryTime string) string {
	if deliveryDate == "" && deliveryTime == "" {
		return "pending"
	}
	return "assigned"
}

// Create delivery object
func deliveryObject(o r.Order, deliveryStatus string) r.Delivery {
	return r.Delivery{
		OrderID:          o.OrderID,
		DeliveryPersonID: 1,
		DeliveryStatus:   deliveryStatus,
		DeliveryDate:     "", // Keep it empty for now until the delivery is actually made
		DeliveryTime:     "", // Keep it empty for now until the delivery is actually made
	}
}

func (h *Handler) UpdateDelivery(c *gin.Context) {
	var delivery r.Delivery
	if err := c.ShouldBindJSON(&delivery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update delivery status in Delivery table
	if err := h.repo.UpdateDelivery(delivery.DeliveryID, delivery.DeliveryStatus, delivery.DeliveryDate, delivery.DeliveryTime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send delivery message
	h.producer.SendDeliveryMessage(delivery)

	c.JSON(http.StatusOK, gin.H{"message": "Delivery status updated successfully"})
}
