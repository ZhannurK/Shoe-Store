package client

import (
	pb "api-gateway/proto/inventory"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var inventoryClient pb.InventoryServiceClient

func InitInventoryGRPCClient() {
	conn, err := grpc.NewClient("inventory-service:5052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("inventory gRPC dial: %v", err)
	}
	inventoryClient = pb.NewInventoryServiceClient(conn)
}

func GRPCGetSneakers(role string, page, limit int32) (*pb.GetSneakersResponse, error) {
	req := &pb.GetSneakersRequest{
		Role:  role,
		Page:  page,
		Limit: limit,
	}
	return inventoryClient.GetSneakers(context.Background(), req)
}

func GRPCCreateSneaker(role, brand, model string, price int32, color string) (*pb.SneakerResponse, error) {
	req := &pb.CreateSneakerRequest{
		Role:  role,
		Brand: brand,
		Model: model,
		Price: price,
		Color: color,
	}
	return inventoryClient.CreateSneaker(context.Background(), req)
}

func GRPCEditSneaker(role, id, brand, model string, price int32, color string) (*pb.SneakerResponse, error) {
	req := &pb.EditSneakerRequest{
		Role:  role,
		Id:    id,
		Brand: brand,
		Model: model,
		Price: price,
		Color: color,
	}
	return inventoryClient.EditSneaker(context.Background(), req)
}

func GRPCRemoveSneaker(role, id string) (*pb.RemoveSneakerResponse, error) {
	req := &pb.RemoveSneakerRequest{
		Role: role,
		Id:   id,
	}
	return inventoryClient.RemoveSneaker(context.Background(), req)
}

func GRPCGetPublicSneakers(page, limit int32) (*pb.GetPublicSneakersResponse, error) {
	req := &pb.GetPublicSneakersRequest{
		Page:  page,
		Limit: limit,
	}
	return inventoryClient.GetPublicSneakers(context.Background(), req)
}
