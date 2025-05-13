package repository

import (
	"context"

	"github.com/shoe-store/inventory-service/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	CreateSneaker(ctx context.Context, sneaker *models.Sneaker) (*models.Sneaker, error)
	GetSneakers(ctx context.Context, page, limit int32) ([]*models.Sneaker, int64, error)
	GetPublicSneakers(ctx context.Context, page, limit int32) ([]*models.Sneaker, int64, error)
	UpdateSneaker(ctx context.Context, id primitive.ObjectID, sneaker *models.Sneaker) error
	DeleteSneaker(ctx context.Context, id primitive.ObjectID) error
	GetSneakerByID(ctx context.Context, id primitive.ObjectID) (*models.Sneaker, error)
	UpdateSneakerStock(ctx context.Context, id primitive.ObjectID, newStock int) error
}
