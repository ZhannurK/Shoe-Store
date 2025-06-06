package client

import (
	proto "api-gateway/proto/user"
	"context"
	"log"

	"google.golang.org/grpc"
)

var authClient proto.AuthServiceClient

func InitAuthGRPCClient() {
	conn, err := grpc.Dial("auth-service:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to auth-service gRPC: %v", err)
	}
	authClient = proto.NewAuthServiceClient(conn)
}

func GRPCSignUp(email, password, name string) (*proto.AuthResponse, error) {
	return authClient.SignUp(context.Background(), &proto.SignUpRequest{
		Email:    email,
		Password: password,
		Name:     name,
	})
}

func GRPCLogin(email, password string) (*proto.AuthResponse, error) {
	return authClient.Login(context.Background(), &proto.LoginRequest{
		Email:    email,
		Password: password,
	})
}

func GRPCConfirmEmail(token string) error {
	_, err := authClient.ConfirmEmail(
		context.Background(),
		&proto.ConfirmEmailRequest{Token: token},
	)
	return err
}

func GRPCChangePassword(userID, oldPwd, newPwd string) error {
	_, err := authClient.ChangePassword(
		context.Background(),
		&proto.ChangePasswordRequest{
			UserId:      userID,
			OldPassword: oldPwd,
			NewPassword: newPwd,
		},
	)
	return err
}
