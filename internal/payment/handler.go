package payment

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

type OrderCreatedMessage struct {
	OrderID     string `json:"order_id"`
	TotalAmount string `json:"total_amount"`
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
		log.Printf("Processing message: %v", string(msg.Value))
		var order OrderCreatedMessage
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Printf("Failed to unmarshal OrderCreated event: %v", err)
			continue
		}

		// Process payment
		if err := h.ProcessPayment(context.Background(), order); err != nil {
			log.Printf("Payment processing failed: %v", err)
		}

		session.MarkMessage(msg, "")
	}
	return nil
}

// ProcessPayment handles payment logic
func (h *Handler) ProcessPayment(ctx context.Context, order OrderCreatedMessage) error {

	// Convert total amount to float64
	totalAmount, err := strconv.ParseFloat(order.TotalAmount, 64)
	if err != nil {
		log.Printf("Failed to convert total amount to float64: %v", err)
	}
	paymentStatus := mockPayment(totalAmount)

	orderID, err := strconv.Atoi(order.OrderID)
	if err != nil {
		log.Printf("Failed to convert order ID to int: %v", err)
	}

	payment := repository.Payment{
		PaymentID:     generatePaymentId(),
		OrderID:       orderID,
		PaymentStatus: paymentStatus,
		PaymentDate:   time.Now(),
	}

	// Save payment to database
	if err := h.repo.AddPayment(payment); err != nil {
		return err
	}

	// Update status of order and not payment
	var orderStatus string
	if paymentStatus == "success" {
		orderStatus = "paid"
	} else {
		orderStatus = "unpaid"
	}
	if err := h.repo.UpdateStatus(orderID, orderStatus); err != nil {
		log.Printf("Failed to update order status: %v", err)
		return err
	}

	paymentData := map[string]string{
		"order_id":       order.OrderID,
		"payment_id":     strconv.Itoa(payment.PaymentID),
		"payment_status": paymentStatus,
		"payment_date":   payment.PaymentDate.Format(time.RFC3339),
	}

	h.producer.SendPaymentMessage(paymentData)

	return nil
}

// generatePaymentId generates a new payment ID.
func generatePaymentId() int {
	return int(time.Now().Unix())
}

// mockPayment processes a payment. This function is not implemented yet so we assume that the payment is a success.
func mockPayment(amount float64) string {
	// Use amount to determine if payment is successful
	if amount < 0 {
		return "failed"
	}
	return "success"
}
