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
	OrderID      int `gorm:"primaryKey"`
	CustomerID   int
	OrderDate    time.Time
	TotalAmount  float64
	OrderStatus  string // e.g. "created", "paid", "delivered"
	DeliveryDate string
	DeliveryTime string
}

// Item represents a list of items in an order.
type Item struct {
	Name      string
	UnitPrice float64
	Quantity  int
	OrderID   int
}

type Payment struct {
	PaymentID     int `gorm:"primaryKey"`
	OrderID       int
	PaymentStatus string // "pending", "success", "failed"
	PaymentDate   time.Time
}

type Delivery struct {
	DeliveryID       int `gorm:"primaryKey"`
	OrderID          int
	DeliveryPersonID int
	DeliveryStatus   string // "pending", "assigned", "delivered", "failed"
}

// DBRepository provides access to the Order storage.
type DBRepository struct {
	db *gorm.DB
}

// Repository provides access to the order storage.
type Repository interface {
	// Adds
	AddOrder(order Order) error
	AddItem(item Item) error
	AddPayment(payment Payment) error
	AddDelivery(delivery Delivery) error

	// Gets
	GetTotalAmount(orderID int) (float64, error)

	// Updates
	UpdateStatus(orderID int, orderStatus string) error
}

// AddOrder adds a new order to the database.
func (receiver *DBRepository) AddOrder(order Order) error {
	return receiver.db.Create(&order).Error
}

// AddItem adds a new item to the database.
func (receiver *DBRepository) AddItem(item Item) error {
	return receiver.db.Create(&item).Error
}

// AddPayment adds a new payment to the database.
func (receiver *DBRepository) AddPayment(payment Payment) error {
	return receiver.db.Create(&payment).Error
}

// AddDelivery adds a new delivery to the database.
func (receiver *DBRepository) AddDelivery(delivery Delivery) error {
	return receiver.db.Create(&delivery).Error
}

// GetTotalAmount gets the total amount of the order.
func (receiver *DBRepository) GetTotalAmount(orderID int) (float64, error) {
	var total float64
	if err := receiver.db.Raw("SELECT SUM(unit_price * quantity) FROM items WHERE order_id = ?", orderID).Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// UpdateStatus updates the status of the payment in the order table.
func (r *DBRepository) UpdateStatus(orderID int, orderStatus string) error {
	return r.db.Model(&Order{}).Where("order_id = ?", orderID).Update("order_status", orderStatus).Error
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
