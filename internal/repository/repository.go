package repository

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	DatabaseURL = "DATABASE_URL"
)

// Order represents an order in the system.
type Order struct {
	OrderID      int       `json:"OrderID" gorm:"primaryKey;autoIncrement"`
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
	OrderID   int     `json:"OrderID"`
}

type Payment struct {
	PaymentID     int       `json:"PaymentID" gorm:"primaryKey;autoIncrement"`
	OrderID       int       `json:"OrderID"`
	PaymentStatus string    `json:"PaymentStatus"` // "pending", "success", "failed"
	PaymentDate   time.Time `json:"PaymentDate"`
}

type Delivery struct {
	DeliveryID       int    `json:"DeliveryID" gorm:"primaryKey;autoIncrement"`
	OrderID          int    `json:"OrderID"`
	DeliveryPersonID int    `json:"DeliveryPersonID"`
	DeliveryStatus   string `json:"DeliveryStatus"` // "pending", "assigned", "delivered", "failed"
	DeliveryDate     string `json:"DeliveryDate"`   // The actual delivery date
	DeliveryTime     string `json:"DeliveryTime"`   // The actual delivery time
}

type Notification struct {
	NotificationID int    `json:"NotificationID" gorm:"primaryKey;autoIncrement"`
	OrderID        int    `json:"OrderID"`
	Message        string `json:"Message"`
}

// DBRepository provides access to the Order storage.
type DBRepository struct {
	db *gorm.DB
}

// Repository provides access to the order storage.
type Repository interface {
	// Adds
	AddOrder(order Order) (int, error)
	AddItem(item Item) error
	AddPayment(payment Payment) (int, error)
	AddDelivery(delivery Delivery) (int, error)
	AddNotification(notification Notification) error

	// Updates
	UpdateStatus(orderID int, orderStatus string) error
	UpdateDelivery(deliveryID int, deliveryStatus string, deliveryDate string, deliveryTime string) error
}

// AddOrder adds a new order to the database and returns order ID.
func (receiver *DBRepository) AddOrder(order Order) (int, error) {
	if err := receiver.db.Create(&order).Error; err != nil {
		return 0, err
	}
	return order.OrderID, nil
}

// AddItem adds a new item to the database.
func (receiver *DBRepository) AddItem(item Item) error {
	return receiver.db.Create(&item).Error
}

// AddPayment adds a new payment to the database.
func (receiver *DBRepository) AddPayment(payment Payment) (int, error) {
	if err := receiver.db.Create(&payment).Error; err != nil {
		return 0, err
	}
	return payment.PaymentID, nil
}

// AddDelivery adds a new delivery to the database.
func (receiver *DBRepository) AddDelivery(delivery Delivery) (int, error) {
	if err := receiver.db.Create(&delivery).Error; err != nil {
		return 0, err
	}
	return delivery.DeliveryID, nil
}

// AddNotification adds a new notification to the database.
func (receiver *DBRepository) AddNotification(notification Notification) error {
	return receiver.db.Create(&notification).Error
}

// UpdateStatus updates the status of the payment in the order table.
func (r *DBRepository) UpdateStatus(orderID int, orderStatus string) error {
	return r.db.Model(&Order{}).Where("order_id = ?", orderID).Update("order_status", orderStatus).Error
}

// UpdateDelivery updates the status of the delivery in the delivery table.
func (r *DBRepository) UpdateDelivery(deliveryID int, deliveryStatus string, deliveryDate string, deliveryTime string) error {
	return r.db.Model(&Delivery{}).Where("delivery_id = ?", deliveryID).Update("delivery_status", deliveryStatus).Update("delivery_date", deliveryDate).Update("delivery_time", deliveryTime).Error
}

// Connect to the database
func ConnectToDatabase() *gorm.DB {
	databaseURL := os.Getenv(DatabaseURL)
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	return db
}

// Automatically migrate the schema
func AutoMigrate(db *gorm.DB, models ...interface{}) {
	if err := db.AutoMigrate(models...); err != nil {
		log.Fatalf("Error migrating the schema: %v", err)
	}
}

// NewDBRepository creates a new Repository backed by a SQL database.
func NewDBRepository(db *gorm.DB) *DBRepository {
	return &DBRepository{db: db}
}
