package natsadapter

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
)

type OrderCreatedEvent struct {
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func PublishOrderCreated(nc *nats.Conn, event OrderCreatedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = nc.Publish("order.created", data)
	if err != nil {
		return err
	}
	log.Println("[NATS] Published event:", event)
	return nil
}
