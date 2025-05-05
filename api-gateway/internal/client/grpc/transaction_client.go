package client

import (
	pb "api-gateway/proto/transaction"
	"context"
	"google.golang.org/grpc"
	"log"
)

var txClient pb.TransactionServiceClient

func InitTransactionGRPCClient() {
	conn, err := grpc.Dial("transaction:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("tx gRPC dial: %v", err)
	}
	txClient = pb.NewTransactionServiceClient(conn)
}

func GRPCCreateTransaction(userID string, items []*pb.CartItem, totalAmount float64) (*pb.Transaction, error) {
	req := &pb.CreateTransactionRequest{
		Transaction: &pb.Transaction{
			UserId:      userID,
			CartItems:   items,
			TotalAmount: totalAmount,
		},
	}

	resp, err := txClient.CreateTransaction(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return resp.Transaction, nil
}

func GRPCGetTransaction(id string) (*pb.Transaction, error) {
	resp, err := txClient.GetTransaction(context.Background(),
		&pb.GetTransactionRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return resp.Transaction, nil
}

func GRPCUpdateStatus(id, status string) error {
	_, err := txClient.UpdateTransactionStatus(context.Background(),
		&pb.UpdateTransactionStatusRequest{Id: id, Status: status})
	return err
}

func GRPCDeleteTransaction(id string) error {
	_, err := txClient.DeleteTransaction(context.Background(),
		&pb.DeleteTransactionRequest{Id: id})
	return err
}
