package natsadapter

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/shoe-store/inventory-service/internal/service"
)

type OrderCreatedEvent struct {
	OrderID   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func SubscribeToOrderCreated(nc *nats.Conn, inventoryService *service.InventoryService) {
	_, err := nc.Subscribe("order.created", func(msg *nats.Msg) {
		var event OrderCreatedEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Println("[NATS] Error unmarshalling:", err)
			return
		}
		log.Println("[NATS] Received event:", event)

		err := inventoryService.DecreaseStock(event.ProductID, event.Quantity)
		if err != nil {
			log.Println("[Inventory] Failed to decrease stock:", err)
		} else {
			log.Println("[Inventory] Stock decreased for product", event.ProductID)
		}
	})

	if err != nil {
		log.Fatal("[NATS] Subscription failed:", err)
	}
}
