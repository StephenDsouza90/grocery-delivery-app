# Grocery Delivery App

## Overview
The **Grocery Delivery App**, written in **Go**, is a simple **microservices** application through which customers grocery orders can easily be processed and delivered to them on a given schedule.

The application has the following services:
- Order Service
- Payment Service
- Delivery Service
- Notification Service (TODO)

![design](assets/images/grocery-delivery-app.png)

It uses a **Postgres** database to store records and **Kafka** along with **Zookeeper** to keep track and process orders between services.

## Table of Contents
- Installation
- Structure overview
- Kafka design
- How to run locally

## Installation

To run the application, **Docker** is required since each service has  `Dockerfile` and a `Postgres`, `Kafka` and `Zookeeper` Docker images are used. The `docker-compose.yml` handles pulling the images of the tools and also builds the services.

## Kafka Design

This part explains how each service talks to each other through `producers` and `consumers`.

![kafka](assets/images/kafka-design.png)

### OrderService

As soon as an order is made by a customer (assuming through a UI), the UI sends a POST request to the `/orders` REST endpoint. The `OrderService` generates an `OrderID` and adds the grocery items to the `item` table and order to the `order` table along with the `OrderID`. The `OrderService` also `produce` a message to the `OrderCreatedTopic`.

### PaymentService

The `PaymentService` subscribes to the `OrderCreatedTopic` to `consume` new messages. The `PaymentService` is responsible for processing the payment, adding the record to the database and `produce` a new message to the `PaymentStatusTopic`.

### DeliveryService

The `DeliveryService` subscribes to the `OrderCreatedTopic` and the `PaymentStatusTopic` to `consume` new messages. From the `PaymentStatusTopic`, it checks if the payment is successful or unsuccessful and from the `OrderCreatedTopic`, it gets the delivery date and time which the customer selected. The `DeliveryService` is responsible for assigning the order to a person for delivery, as well as adding the record to the database and `produce` a message to the `DeliveryStatusTopic`.

### NotificationService

TODO 

## How to run locally

To run the app locally, use the following steps:

Use the `docker-compose.yml` to set up `postgres`, `kafka` and `zookeeper` by running the following command:

```
docker-compose up -d
```

To restart everything:

```
docker-compose up --build -d
```

To stop the services, use the following command:
```
docker-compose down
```

To check if the containers are running, running the following command:
```
docker ps
```

In a new terminal, a CURL request can be sent to the `/orders` REST endpoint as below:

```curl
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": 1,
    "items": {
        "item1": {
            "name": "item1", 
            "quantity": 2,
            "unitPrice": 3
        },
        "item2": {
            "name": "item2",
            "quantity": 3,
            "unitPrice": 5
        }
    },
    "delivery_date": "2025-03-07",
    "delivery_time": "14:00"
  }'
```

To check if the records are saved to the database, connect to the docker container by running the following command:

```
docker exec -it grocery-delivery-app-postgres-1 psql -U postgres -d grocerydelivery
```

To find all tables run:
```
\dt
```

To query the tables, the below SQL statement can be executed:
- `SELECT * FROM items;`
- `SELECT * FROM orders;`
- `SELECT * FROM payments;`
- `SELECT * FROM deliveries;`

To see the messages in Kafka, connect to the Kafka container:
```
docker exec -it kafka /bin/bash
```

Once inside the Kafka container, use the kafka-console-consumer to see messages from a specific topic.
- `kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic OrderCreated --from-beginning`
- `kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic PaymentStatus --from-beginning`
- `kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic DeliveryStatus --from-beginning`

To restart all volumes, run:
```
docker volume rm $(docker volume ls -q)
```

Remove volumes to reset data

```
docker-compose down -v 
```


## TODO:
- Write test cases
- Dead Letter Queue for failed payments
- Central logging
- Use GraphQL instead of REST
- More kafka standards
- More features
- UI to send many (random) requests
- Analysis of each step
- Review images


list all Kafka topics
docker exec kafka kafka-topics.sh --list --bootstrap-server kafka:9092

- Delivery person delivers order and updates status `/delivery-update` -> DeliveryStatus changes to completed -> OrderStatus changes to delivered
- Notification services -> Subscribes to order, payment and delivery and either produces messages which are sent to the customer and saved in the database