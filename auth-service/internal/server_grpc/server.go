package server_grpc

import (
	"auth-service/internal/usecase"
	"auth-service/proto"
	"context"
)

type AuthGRPCServer struct {
	proto.UnimplementedAuthServiceServer
	UC usecase.AuthUseCase
}

func NewAuthGRPCServer(uc usecase.AuthUseCase) *AuthGRPCServer {
	return &AuthGRPCServer{UC: uc}
}

func (s *AuthGRPCServer) SignUp(ctx context.Context, req *proto.SignUpRequest) (*proto.AuthResponse, error) {
	token, err := s.UC.SignUp(ctx, req.Email, req.Password, req.Name)
	if err != nil {
		return nil, err
	}
	return &proto.AuthResponse{Token: token}, nil
}

func (s *AuthGRPCServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.AuthResponse, error) {
	token, _, err := s.UC.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &proto.AuthResponse{Token: token}, nil
}

func (s *AuthGRPCServer) ConfirmEmail(ctx context.Context, req *proto.ConfirmEmailRequest) (*proto.Empty, error) {
	err := s.UC.ConfirmEmail(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}

func (s *AuthGRPCServer) ChangePassword(ctx context.Context, req *proto.ChangePasswordRequest) (*proto.Empty, error) {
	err := s.UC.ChangePassword(ctx, req.UserId, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, err
	}
	return &proto.Empty{}, nil
}
