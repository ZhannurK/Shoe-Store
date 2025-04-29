package service

import (
	"context"
	"errors"

	"github.com/shoe-store/inventory-service/internal/models"
	"github.com/shoe-store/inventory-service/internal/repository"
	"github.com/shoe-store/inventory-service/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type InventoryService struct {
	proto.UnimplementedInventoryServiceServer
	repo repository.Repository
}

func NewInventoryService(repo repository.Repository) *InventoryService {
	return &InventoryService{
		repo: repo,
	}
}

func (s *InventoryService) GetSneakers(ctx context.Context, req *proto.GetSneakersRequest) (*proto.GetSneakersResponse, error) {
	if req.Role != RoleAdmin {
		return nil, errors.New("permission denied: only admins can access this endpoint")
	}

	sneakers, total, err := s.repo.GetSneakers(ctx, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	protoSneakers := make([]*proto.Sneaker, len(sneakers))
	for i, s := range sneakers {
		protoSneakers[i] = &proto.Sneaker{
			Id:    s.ID.Hex(),
			Brand: s.Brand,
			Model: s.Model,
			Price: int32(s.Price),
			Color: s.Color,
		}
	}

	return &proto.GetSneakersResponse{
		Sneakers: protoSneakers,
		Total:    int32(total),
	}, nil
}

func (s *InventoryService) CreateSneaker(ctx context.Context, req *proto.CreateSneakerRequest) (*proto.SneakerResponse, error) {
	if req.Role != RoleAdmin {
		return nil, errors.New("permission denied: only admins can access this endpoint")
	}

	if req.Brand == "" || req.Model == "" || req.Price <= 0 || req.Color == "" {
		return nil, errors.New("invalid input: all fields are required")
	}

	sneaker := &models.Sneaker{
		Brand: req.Brand,
		Model: req.Model,
		Price: int(req.Price),
		Color: req.Color,
	}

	createdSneaker, err := s.repo.CreateSneaker(ctx, sneaker)
	if err != nil {
		return nil, err
	}

	return &proto.SneakerResponse{
		Sneaker: &proto.Sneaker{
			Id:    createdSneaker.ID.Hex(),
			Brand: createdSneaker.Brand,
			Model: createdSneaker.Model,
			Price: int32(createdSneaker.Price),
			Color: createdSneaker.Color,
		},
	}, nil
}

func (s *InventoryService) EditSneaker(ctx context.Context, req *proto.EditSneakerRequest) (*proto.SneakerResponse, error) {
	if req.Role != RoleAdmin {
		return nil, errors.New("permission denied: only admins can access this endpoint")
	}

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, errors.New("invalid sneaker ID")
	}

	sneaker := &models.Sneaker{
		Brand: req.Brand,
		Model: req.Model,
		Price: int(req.Price),
		Color: req.Color,
	}

	err = s.repo.UpdateSneaker(ctx, id, sneaker)
	if err != nil {
		return nil, err
	}

	return &proto.SneakerResponse{
		Sneaker: &proto.Sneaker{
			Id:    id.Hex(),
			Brand: sneaker.Brand,
			Model: sneaker.Model,
			Price: int32(sneaker.Price),
			Color: sneaker.Color,
		},
	}, nil
}

func (s *InventoryService) RemoveSneaker(ctx context.Context, req *proto.RemoveSneakerRequest) (*proto.RemoveSneakerResponse, error) {
	if req.Role != RoleAdmin {
		return nil, errors.New("permission denied: only admins can access this endpoint")
	}

	id, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, errors.New("invalid sneaker ID")
	}

	err = s.repo.DeleteSneaker(ctx, id)
	if err != nil {
		return nil, err
	}

	return &proto.RemoveSneakerResponse{
		Success: true,
	}, nil
}

func (s *InventoryService) GetPublicSneakers(ctx context.Context, req *proto.GetPublicSneakersRequest) (*proto.GetPublicSneakersResponse, error) {
	sneakers, total, err := s.repo.GetPublicSneakers(ctx, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	protoSneakers := make([]*proto.Sneaker, len(sneakers))
	for i, s := range sneakers {
		protoSneakers[i] = &proto.Sneaker{
			Id:    s.ID.Hex(),
			Brand: s.Brand,
			Model: s.Model,
			Price: int32(s.Price),
			Color: s.Color,
		}
	}

	return &proto.GetPublicSneakersResponse{
		Sneakers: protoSneakers,
		Total:    int32(total),
	}, nil
}
