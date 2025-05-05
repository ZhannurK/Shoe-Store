package usecase

import (
	"context"
	"github.com/nats-io/nats.go"
	"log"

	"transaction-service/internal/domain"
	"transaction-service/internal/natsadapter"
	"transaction-service/internal/repositories"
)

type TransactionUseCase struct {
	Repo     repositories.TransactionRepository
	NatsConn *nats.Conn
}

func NewTransactionUseCase(repo repositories.TransactionRepository, nc *nats.Conn) *TransactionUseCase {
	return &TransactionUseCase{
		Repo:     repo,
		NatsConn: nc,
	}
}

func (u *TransactionUseCase) Create(ctx context.Context, tx *domain.Transaction) error {
	return u.Repo.Create(ctx, tx)
}

func (u *TransactionUseCase) UpdateStatus(ctx context.Context, transactionID string, status domain.TransactionStatus) error {
	err := u.Repo.UpdateStatus(ctx, transactionID, status)
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
	return u.Repo.GetByID(ctx, transactionID)
}

func (u *TransactionUseCase) DeleteTransaction(ctx context.Context, id string) error {
	return u.Repo.Delete(ctx, id)
}
