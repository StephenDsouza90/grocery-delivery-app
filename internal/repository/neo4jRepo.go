package repository

import (
	"context"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// DBRepository provides access to the Order storage.
type neo4jDBRepository struct {
	driver neo4j.DriverWithContext
}

// ConnectToDatabase connects to the Neo4j database.
func ConnectToNeo4jDatabase() neo4j.DriverWithContext {
	databaseURL := os.Getenv(DatabaseURL)
	driver, err := neo4j.NewDriverWithContext(databaseURL, neo4j.BasicAuth("neo4j", "password", ""))
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	return driver
}

// NewDBRepository creates a new Repository backed by a Neo4j database.
func NewNeo4jDBRepository(driver neo4j.DriverWithContext) *neo4jDBRepository {
	return &neo4jDBRepository{driver: driver}
}

// Set up constraints and indexes
func SetupConstraints(driver neo4j.DriverWithContext) {
	session := driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.TODO())

	constraints := []string{
		"CREATE CONSTRAINT Order_OrderID_unique FOR (o:Order) REQUIRE o.OrderID IS UNIQUE;",
		"CREATE CONSTRAINT Payment_PaymentID_unique FOR (p:Payment) REQUIRE p.PaymentID IS UNIQUE",
		"CREATE CONSTRAINT Delivery_DeliveryID_unique FOR (d:Delivery) REQUIRE d.DeliveryID IS UNIQUE",
		"CREATE CONSTRAINT Notification_NotificationID_unique FOR (n:Notification) REQUIRE n.NotificationID IS UNIQUE",
	}

	for _, constraint := range constraints {
		_, err := session.Run(context.TODO(), constraint, nil)
		if err != nil {
			log.Fatalf("Error setting up constraints: %v", err)
		}
	}
}

// AddOrder adds a new order to the database and returns order ID.
func (r *neo4jDBRepository) AddOrder(order Order) error {
	session := r.driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.TODO())

	_, err := session.Run(
		context.TODO(),
		"CREATE (o:Order {OrderID: $OrderID, CustomerID: $CustomerID, OrderDate: $OrderDate, TotalAmount: $TotalAmount, OrderStatus: $OrderStatus, DeliveryDate: $DeliveryDate, DeliveryTime: $DeliveryTime}) RETURN o.OrderID",
		map[string]interface{}{
			"OrderID":      order.OrderID,
			"CustomerID":   order.CustomerID,
			"OrderDate":    order.OrderDate,
			"TotalAmount":  order.TotalAmount,
			"OrderStatus":  order.OrderStatus,
			"DeliveryDate": order.DeliveryDate,
			"DeliveryTime": order.DeliveryTime,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// AddItem adds a new item to the database.
func (r *neo4jDBRepository) AddItem(item Item) error {
	session := r.driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.TODO())

	_, err := session.Run(
		context.TODO(),
		"CREATE (i:Item {Name: $Name, UnitPrice: $UnitPrice, Quantity: $Quantity, OrderID: $OrderID})",
		map[string]interface{}{
			"Name":      item.Name,
			"UnitPrice": item.UnitPrice,
			"Quantity":  int64(item.Quantity),
			"OrderID":   item.OrderID,
		},
	)
	return err
}

// AddPayment adds a new payment to the database.
func (r *neo4jDBRepository) AddPayment(payment Payment) error {
	session := r.driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.TODO())

	_, err := session.Run(
		context.TODO(),
		"CREATE (p:Payment {PaymentID: $PaymentID, OrderID: $OrderID, PaymentStatus: $PaymentStatus, PaymentDate: $PaymentDate}) RETURN p.PaymentID",
		map[string]interface{}{
			"PaymentID":     payment.PaymentID,
			"OrderID":       payment.OrderID,
			"PaymentStatus": payment.PaymentStatus,
			"PaymentDate":   payment.PaymentDate,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

// AddDelivery adds a new delivery to the database.
func (r *neo4jDBRepository) AddDelivery(delivery Delivery) error {
	session := r.driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.TODO())

	_, err := session.Run(
		context.TODO(),
		"CREATE (d:Delivery {DeliveryID: $DeliveryID, OrderID: $OrderID, DeliveryPersonID: $DeliveryPersonID, DeliveryStatus: $DeliveryStatus, DeliveryDate: $DeliveryDate, DeliveryTime: $DeliveryTime}) RETURN d.DeliveryID",
		map[string]interface{}{
			"DeliveryID":       delivery.DeliveryID,
			"OrderID":          delivery.OrderID,
			"DeliveryPersonID": delivery.DeliveryPersonID,
			"DeliveryStatus":   delivery.DeliveryStatus,
			"DeliveryDate":     delivery.DeliveryDate,
			"DeliveryTime":     delivery.DeliveryTime,
		},
	)

	if err != nil {
		return err
	}
	return nil
}

// AddNotification adds a new notification to the database.
func (r *neo4jDBRepository) AddNotification(notification Notification) error {
	session := r.driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.TODO())

	_, err := session.Run(
		context.TODO(),
		"CREATE (n:Notification {NotificationID: $NotificationID, OrderID: $OrderID, Message: $Message})",
		map[string]interface{}{
			"NotificationID": notification.NotificationID,
			"OrderID":        notification.OrderID,
			"Message":        notification.Message,
		},
	)
	return err
}

func (r *neo4jDBRepository) UpdateStatus(orderID string, orderStatus string) error {
	session := r.driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.TODO())

	_, err := session.Run(
		context.TODO(),
		"MATCH (o:Order {OrderID: $OrderID}) SET o.OrderStatus = $OrderStatus",
		map[string]interface{}{
			"OrderID":     orderID,
			"OrderStatus": orderStatus,
		},
	)
	return err
}

func (r *neo4jDBRepository) UpdateDelivery(deliveryID string, deliveryStatus string, deliveryDate string, deliveryTime string) error {
	session := r.driver.NewSession(context.TODO(), neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(context.TODO())

	_, err := session.Run(
		context.TODO(),
		"MATCH (d:Delivery {DeliveryID: $DeliveryID}) SET d.DeliveryStatus = $DeliveryStatus, d.DeliveryDate = $DeliveryDate, d.DeliveryTime = $DeliveryTime",
		map[string]interface{}{
			"DeliveryID":     deliveryID,
			"DeliveryStatus": deliveryStatus,
			"DeliveryDate":   deliveryDate,
			"DeliveryTime":   deliveryTime,
		},
	)
	return err
}
