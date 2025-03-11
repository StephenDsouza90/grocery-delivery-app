package handler

import (
	"net/http"
	"time"

	k "github.com/StephenDsouza90/grocery-delivery-app/internal/kafka"
	r "github.com/StephenDsouza90/grocery-delivery-app/internal/repository"
	u "github.com/StephenDsouza90/grocery-delivery-app/internal/utils"
	"github.com/gin-gonic/gin"
)

// Handler provides HTTP handlers for interacting with orders.
type Handler struct {
	repo     r.Repository
	producer *k.Producer
	consumer *k.Consumer
}

// OrderRequestBody represents a request to create a new order.
type OrderRequestBody struct {
	CustomerID   int               `json:"CustomerID"`
	Items        map[string]r.Item `json:"Items"`
	DeliveryDate string            `json:"DeliveryDate"` // Format: "2006-01-02"
	DeliveryTime string            `json:"DeliveryTime"` // Format: "13:00 to 14:00"
}

// NewHandler creates a new Handler.
func NewHandler(repo r.Repository, producer *k.Producer, consumer *k.Consumer) *Handler {
	return &Handler{
		repo:     repo,
		producer: producer,
		consumer: consumer,
	}
}

// CreateOrder creates a new order.
func (h *Handler) CreateOrder(c *gin.Context) {
	var rb OrderRequestBody
	if err := c.ShouldBindJSON(&rb); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	totalAmount := getTotalAmount(rb.Items)

	// Create new order object and save it in the database
	newOrder := orderObject(rb, totalAmount)
	err := h.repo.AddOrder(newOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create and save items in the order
	for _, item := range rb.Items {
		item.OrderID = newOrder.OrderID
		if err := h.repo.AddItem(item); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Send a message to Kafka
	h.producer.SendOrderMessage(newOrder)

	c.JSON(http.StatusCreated, newOrder)
}

// getTotalAmount calculates the total amount of the order.
func getTotalAmount(items map[string]r.Item) float64 {
	var total float64
	for _, item := range items {
		total += item.UnitPrice * float64(item.Quantity)
	}
	return total
}

// orderObject creates a new order object.
func orderObject(rb OrderRequestBody, totalAmount float64) r.Order {
	return r.Order{
		OrderID:      u.GenerateUUID(),
		CustomerID:   rb.CustomerID,
		OrderDate:    time.Now(),
		TotalAmount:  totalAmount,
		OrderStatus:  "created",
		DeliveryDate: rb.DeliveryDate,
		DeliveryTime: rb.DeliveryTime,
	}
}
