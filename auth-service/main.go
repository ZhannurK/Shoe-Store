package main

import (
	"auth-service/internal/cache"
	"auth-service/internal/handler"
	"auth-service/internal/repositories"
	"auth-service/internal/server_grpc"
	"auth-service/internal/usecase"
	"auth-service/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var requestCount = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "shoestore_requests_total",
		Help: "Total number of HTTP requests received",
	},
)

func main() {

	client := repositories.ConnectDB()

	jwtKey := []byte(os.Getenv("JWTSECRET"))

	userRepo := repositories.NewUserMongoRepo(client, "db", "users")

	redisCache := cache.NewRedisCache("redis:6379")
	authUC := usecase.NewAuthUseCase(userRepo, jwtKey, redisCache)

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

		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				requestCount.Inc()
				next.ServeHTTP(w, req)
			})
		})

		prometheus.MustRegister(requestCount)
		r.Handle("/metrics", promhttp.Handler())

		log.Println("🌐 REST server running on port 8087...")
		if err := http.ListenAndServe(":8087", r); err != nil {
			log.Fatalf("Failed to start REST server: %v", err)
		}
	}()

	select {}
}
