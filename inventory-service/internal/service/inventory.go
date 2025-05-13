package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/shoe-store/inventory-service/internal/cache"
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

	// Инвалидация кеша всех списков
	s.invalidatePublicSneakersCache(ctx)

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

	s.invalidatePublicSneakersCache(ctx)

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

	s.invalidatePublicSneakersCache(ctx)

	return &proto.RemoveSneakerResponse{
		Success: true,
	}, nil
}

func (s *InventoryService) GetPublicSneakers(ctx context.Context, req *proto.GetPublicSneakersRequest) (*proto.GetPublicSneakersResponse, error) {
	cacheKey := fmt.Sprintf("public_sneakers_list:%d:%d", req.Page, req.Limit)

	// 1. Пробуем получить из кеша
	cached, err := cache.Rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		log.Printf("[CACHE HIT] Key: %s", cacheKey)
		var resp proto.GetPublicSneakersResponse
		if err := json.Unmarshal([]byte(cached), &resp); err == nil {
			return &resp, nil
		}
	}
	log.Printf("[CACHE MISS] Key: %s", cacheKey)

	// 2. Если нет в кеше — получаем из БД
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

	result := &proto.GetPublicSneakersResponse{
		Sneakers: protoSneakers,
		Total:    int32(total),
	}

	// 3. Сохраняем в кеш на 5 минут
	bytes, _ := json.Marshal(result)
	cache.Rdb.Set(ctx, cacheKey, bytes, 5*time.Minute)
	log.Printf("[CACHE SET] Key: %s", cacheKey)

	return result, nil
}

func (s *InventoryService) DecreaseStock(productID string, quantity int) error {
	id, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return errors.New("invalid product ID")
	}

	ctx := context.Background()
	sneaker, err := s.repo.GetSneakerByID(ctx, id)
	if err != nil {
		return err
	}

	if sneaker.Stock < quantity {
		return errors.New("not enough stock")
	}

	sneaker.Stock -= quantity

	return s.repo.UpdateSneakerStock(ctx, id, sneaker.Stock)
}

// Удаляем все кеши списков (на случай пагинации)
func (s *InventoryService) invalidatePublicSneakersCache(ctx context.Context) {
	iter := cache.Rdb.Scan(ctx, 0, "public_sneakers_list*", 0).Iterator()
	for iter.Next(ctx) {
		cache.Rdb.Del(ctx, iter.Val())
		log.Printf("[CACHE DEL] Key: %s", iter.Val())
	}
}
