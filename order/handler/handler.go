package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/StephenDsouza90/grocery-delivery-app/internal/kafka"
	"github.com/StephenDsouza90/grocery-delivery-app/internal/repository"
	"github.com/gin-gonic/gin"
)

// Handler provides HTTP handlers for interacting with orders.
type Handler struct {
	repo     repository.Repository
	producer *kafka.Producer
	consumer *kafka.Consumer
}

// CreateOrderRequest represents a request to create a new order.
type CreateOrderRequest struct {
	CustomerID   int                        `json:"customer_id"`
	Items        map[string]repository.Item `json:"items"`
	DeliveryDate string                     `json:"delivery_date"` // Format: "2006-01-02"
	DeliveryTime string                     `json:"delivery_time"` // Format: "13:00 to 14:00"
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

// CreateOrder creates a new order.
func (handler *Handler) CreateOrder(context *gin.Context) {
	var request CreateOrderRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderId := generateOrderNumber()

	// Create items
	for _, item := range request.Items {
		item.OrderID = orderId
		if err := handler.repo.AddItem(item); err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Get total amount
	totalAmount, err := handler.repo.GetTotalAmount(orderId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a new order
	newOrder := repository.Order{
		OrderID:      orderId,
		CustomerID:   request.CustomerID,
		OrderDate:    time.Now(),
		TotalAmount:  totalAmount,
		OrderStatus:  "created",
		DeliveryDate: request.DeliveryDate,
		DeliveryTime: request.DeliveryTime,
	}

	// Save the order
	if err := handler.repo.AddOrder(newOrder); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Prepare the order data
	orderMap := map[string]string{
		"order_id":      strconv.Itoa(newOrder.OrderID),
		"customer_id":   strconv.Itoa(newOrder.CustomerID),
		"order_date":    newOrder.OrderDate.Format(time.RFC3339),
		"total_amount":  fmt.Sprintf("%.2f", newOrder.TotalAmount),
		"order_status":  newOrder.OrderStatus,
		"delivery_date": newOrder.DeliveryDate,
		"delivery_time": newOrder.DeliveryTime,
	}

	// Send a message to Kafka
	handler.producer.SendOrderMessage(orderMap)

	context.JSON(http.StatusCreated, newOrder)
}

// Generate a new order number based on the current time.
func generateOrderNumber() int {
	return int(time.Now().Unix())
}

// Gets total amount of the order
func (request *CreateOrderRequest) Amount() float64 {
	var total float64
	for _, item := range request.Items {
		total += item.UnitPrice * float64(item.Quantity)
	}
	return total
}
