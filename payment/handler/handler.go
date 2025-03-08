package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	k "github.com/StephenDsouza90/grocery-delivery-app/internal/kafka"
	r "github.com/StephenDsouza90/grocery-delivery-app/internal/repository"
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
		var order r.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Failed to unmarshal OrderCreated event: %v", err)
			continue
		}

		if err := h.ProcessPayment(context.Background(), order); err != nil {
			log.Printf("Payment processing failed: %v", err)
		}

		session.MarkMessage(msg, "")
	}
	return nil
}

// ProcessPayment handles payment logic
func (h *Handler) ProcessPayment(c context.Context, o r.Order) error {
	paymentStatus := mockPayment(o.TotalAmount)

	// Create and save payment to database
	newPayment := paymentObject(o, paymentStatus)
	paymentID, err := h.repo.AddPayment(newPayment)
	if err != nil {
		log.Printf("Failed to save payment to database: %v", err)
		return err
	}

	newPayment.PaymentID = paymentID
	h.producer.SendPaymentMessage(newPayment)

	return nil
}

// mockPayment processes a payment. This function is not implemented yet so we assume that the payment is a success.
func mockPayment(amount float64) string {
	// Use amount to determine if payment is successful
	if amount < 0 {
		return "failed"
	}
	return "success"
}

// paymentObject creates a new payment object.
func paymentObject(order r.Order, paymentStatus string) r.Payment {
	return r.Payment{
		OrderID:       order.OrderID,
		PaymentStatus: paymentStatus,
		PaymentDate:   time.Now(),
	}
}
