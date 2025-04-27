package grpc

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"transaction-service/internal/domain"
	"transaction-service/internal/usecase"
	pb "transaction-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedTransactionServiceServer
	uc *usecase.TransactionUseCase
}

func New(uc *usecase.TransactionUseCase) *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterTransactionServiceServer(s, &Server{uc: uc})
	reflection.Register(s)
	return s
}

func (s *Server) CreateTransaction(
	ctx context.Context,
	req *pb.CreateTransactionRequest,
) (*pb.CreateTransactionResponse, error) {

	in := req.GetTransaction()
	if in == nil {
		return nil, status.Error(codes.InvalidArgument, "transaction payload is required")
	}

	tx := &domain.Transaction{
		TransactionID: in.TransactionId,
		UserID:        in.UserId,
		TotalAmount:   in.TotalAmount,
		Status:        domain.TransactionStatus(in.Status),
	}

	tx.CartItems = make([]domain.CartItem, len(in.CartItems))
	for i, c := range in.CartItems {
		oid, err := primitive.ObjectIDFromHex(c.SneakerId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument,
				"invalid sneaker_id %q: %v", c.SneakerId, err)
		}
		tx.CartItems[i] = domain.CartItem{
			SneakerID: oid,
			Quantity:  int(c.Quantity),
		}
	}

	if err := s.uc.Create(ctx, tx); err != nil {
		return nil, err
	}
	return &pb.CreateTransactionResponse{Transaction: tx.ToProto()}, nil
}

func (s *Server) GetTransaction(
	ctx context.Context,
	req *pb.GetTransactionRequest,
) (*pb.GetTransactionResponse, error) {

	tx, err := s.uc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetTransactionResponse{Transaction: tx.ToProto()}, nil
}

func (s *Server) UpdateTransactionStatus(
	ctx context.Context,
	req *pb.UpdateTransactionStatusRequest,
) (*pb.UpdateTransactionStatusResponse, error) {

	if err := s.uc.UpdateStatus(ctx, req.Id,
		domain.TransactionStatus(req.Status)); err != nil {
		return nil, err
	}

	updated, _ := s.uc.GetByID(ctx, req.Id)
	return &pb.UpdateTransactionStatusResponse{
		Transaction: updated.ToProto(),
	}, nil
}

func (s *Server) DeleteTransaction(
	ctx context.Context,
	req *pb.DeleteTransactionRequest,
) (*pb.DeleteTransactionResponse, error) {

	if err := s.uc.DeleteTransaction(ctx, req.Id); err != nil {
		return nil, err
	}
	return &pb.DeleteTransactionResponse{Success: true}, nil
}
