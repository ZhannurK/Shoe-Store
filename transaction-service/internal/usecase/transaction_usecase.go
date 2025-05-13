package usecase

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"time"

	"transaction-service/internal/domain"
	"transaction-service/internal/natsadapter"
	"transaction-service/internal/repositories"
)

type TransactionUseCase struct {
	Repo     repositories.TransactionRepository
	NatsConn *nats.Conn
	Cache    Cache
}

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl time.Duration) error
	Del(key string) error
}

func NewTransactionUseCase(repo repositories.TransactionRepository, cache Cache, nc *nats.Conn) *TransactionUseCase {
	return &TransactionUseCase{Repo: repo, Cache: cache}
}

func (u *TransactionUseCase) Create(ctx context.Context, tx *domain.Transaction) error {
	return u.Repo.Create(ctx, tx)
}

func (u *TransactionUseCase) UpdateStatus(ctx context.Context, transactionID string, status domain.TransactionStatus) error {
	err := u.Repo.UpdateStatus(ctx, transactionID, status)
	if err != nil {
		return err
	}

	err = u.Cache.Del("transaction:" + transactionID)
	if err != nil {
		return err
	}

	if status == domain.StatusPaid {
		tx, err := u.Repo.GetByID(ctx, transactionID)
		if err != nil {
			log.Println("[NATS] Failed to retrieve transaction:", err)
			return nil
		}

		for _, item := range tx.CartItems {
			event := natsadapter.OrderCreatedEvent{
				OrderID:   tx.ID.Hex(),
				ProductID: item.SneakerID.Hex(),
				Quantity:  item.Quantity,
			}
			if len(tx.CartItems) == 0 {
				log.Println("[NATS] No cart items to publish")
				return nil
			}
			if err := natsadapter.PublishOrderCreated(u.NatsConn, event); err != nil {
				log.Println("[NATS] Failed to publish event:", err)
			} else {
				log.Println("[NATS] Published order.created for:", event.OrderID)
			}
		}

	}

	return nil
}

func (u *TransactionUseCase) GetByID(ctx context.Context, transactionID string) (*domain.Transaction, error) {
	cacheKey := "transaction:" + transactionID

	if cached, _ := u.Cache.Get(cacheKey); cached != "" {
		var tx domain.Transaction
		_ = json.Unmarshal([]byte(cached), &tx)
		return &tx, nil
	}

	tx, err := u.Repo.GetByID(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(tx); err == nil {
		err := u.Cache.Set(cacheKey, string(data), 5*time.Minute)
		if err != nil {
			return nil, err
		}
	}

	if cached, _ := u.Cache.Get(cacheKey); cached != "" {
		log.Println("✅ [CACHE HIT] Returning cached value for", cacheKey)
		return tx, nil
	} else {
		log.Println("❌ [CACHE MISS] Fetching from DB and caching", cacheKey)
	}

	return tx, nil
}

func (u *TransactionUseCase) DeleteTransaction(ctx context.Context, id string) error {
	return u.Repo.Delete(ctx, id)
}
