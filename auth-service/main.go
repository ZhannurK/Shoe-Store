package main

import (
	"auth-service/internal/handler"
	"auth-service/internal/repositories"
	"auth-service/internal/server_grpc"
	"auth-service/internal/usecase"
	"auth-service/proto"
	"google.golang.org/grpc"

	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {

	client := repositories.ConnectDB()

	jwtKey := []byte(os.Getenv("JWTSECRET"))

	userRepo := repositories.NewUserMongoRepo(client, "db", "users")

	authUC := usecase.NewAuthUseCase(userRepo, jwtKey)

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen on gRPC: %v", err)
		}

		server := grpc.NewServer()
		proto.RegisterAuthServiceServer(server, server_grpc.NewAuthGRPCServer(authUC))

		log.Println("🚀 gRPC server running on port 50051...")
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	go func() {
		r := mux.NewRouter()
		authHandler := handler.NewAuthHandler(authUC)
		authHandler.RegisterRoutes(r)

		log.Println("🌐 REST server running on port 8087...")
		if err := http.ListenAndServe(":8087", r); err != nil {
			log.Fatalf("Failed to start REST server: %v", err)
		}
	}()

	select {}
}
