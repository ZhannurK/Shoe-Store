package usecase

import (
	"context"
	"transaction-service/internal/domain"
	"transaction-service/internal/repositories"
)

type TransactionUseCase struct {
	Repo repositories.TransactionRepository
}

func NewTransactionUseCase(repo repositories.TransactionRepository) *TransactionUseCase {
	return &TransactionUseCase{Repo: repo}
}

func (u *TransactionUseCase) Create(ctx context.Context, tx *domain.Transaction) error {
	return u.Repo.Create(ctx, tx)
}

func (u *TransactionUseCase) UpdateStatus(ctx context.Context, transactionID string, status domain.TransactionStatus) error {
	return u.Repo.UpdateStatus(ctx, transactionID, status)
}

func (u *TransactionUseCase) GetByID(ctx context.Context, transactionID string) (*domain.Transaction, error) {
	return u.Repo.GetByID(ctx, transactionID)
}

func (u *TransactionUseCase) DeleteTransaction(ctx context.Context, id string) error {
	return u.Repo.Delete(ctx, id)
}
