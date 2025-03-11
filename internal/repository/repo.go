package repository

import (
	"time"
)

const (
	DatabaseURL = "DATABASE_URL"
)

// Order represents an order in the system.
type Order struct {
	OrderID      string    `json:"OrderID"`
	CustomerID   int       `json:"CustomerID"`
	OrderDate    time.Time `json:"OrderDate"`
	TotalAmount  float64   `json:"TotalAmount"`
	OrderStatus  string    `json:"OrderStatus"`  // e.g. "created", "paid", "assigned", "delivered"
	DeliveryDate string    `json:"DeliveryDate"` // The customer's preferred delivery date
	DeliveryTime string    `json:"DeliveryTime"` // The customer's preferred delivery time
}

// Item represents a list of items in an order.
type Item struct {
	Name      string  `json:"Name"`
	UnitPrice float64 `json:"UnitPrice"`
	Quantity  int     `json:"Quantity"`
	OrderID   string  `json:"OrderID"`
}

type Payment struct {
	PaymentID     string    `json:"PaymentID"`
	OrderID       string    `json:"OrderID"`
	PaymentStatus string    `json:"PaymentStatus"` // "pending", "success", "failed"
	PaymentDate   time.Time `json:"PaymentDate"`
}

// Also the delivery request body
type Delivery struct {
	DeliveryID       string `json:"DeliveryID"`
	OrderID          string `json:"OrderID"`
	DeliveryPersonID int    `json:"DeliveryPersonID"`
	DeliveryStatus   string `json:"DeliveryStatus"` // "pending", "assigned", "delivered", "failed"
	DeliveryDate     string `json:"DeliveryDate"`   // The actual delivery date
	DeliveryTime     string `json:"DeliveryTime"`   // The actual delivery time
}

type Notification struct {
	NotificationID string `json:"NotificationID"`
	OrderID        string `json:"OrderID"`
	Message        string `json:"Message"`
}

// Repository provides access to the order storage.
type Repository interface {
	// Adds
	AddOrder(order Order) error
	AddItem(item Item) error
	AddPayment(payment Payment) error
	AddDelivery(delivery Delivery) error
	AddNotification(notification Notification) error

	// Updates
	UpdateStatus(orderID string, orderStatus string) error
	UpdateDelivery(deliveryID string, deliveryStatus string, deliveryDate string, deliveryTime string) error
}
