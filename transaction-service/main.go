package main

import (
	"context"
	"errors"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"transaction-service/internal/cache"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	tgRPC "transaction-service/internal/grpc"
	"transaction-service/internal/handlers"
	"transaction-service/internal/repositories"
	"transaction-service/internal/usecase"
)

var requestCount = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "shoestore_requests_total",
		Help: "Total number of HTTP requests received",
	},
)

func main() {
	//env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system variables")
	}

	mongoURI := os.Getenv("MONGO_CONNECT")
	if mongoURI == "" {
		log.Fatal("MONGO_CONNECT env variable is not set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}

	//repositories / use cases
	mongoClient := repositories.InitMongo(mongoURI)
	txRepo := repositories.NewTransactionRepository(mongoClient, "OnlineStore")
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	redisCache := cache.NewRedisCache("redis:6379")
	txUC := usecase.NewTransactionUseCase(txRepo, redisCache, nc)
	txHandler := handlers.NewTransactionHandler(txUC)

	// ----- gRPC transport --------------------------------------------------
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("gRPC listen: %v", err)
		}
		grpcSrv := tgRPC.New(txUC)
		log.Println("gRPC ⟶ :50051")
		if err := grpcSrv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatalf("gRPC serve: %v", err)
		}
	}()

	// ----- REST transport --------------------------------------------------

	//routers
	r := gin.Default()

	prometheus.MustRegister(requestCount)

	r.Use(func(c *gin.Context) {
		requestCount.Inc()
		c.Next()
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.POST("/transactions", txHandler.CreateTransaction)
	r.GET("/transactions/:id", txHandler.GetTransaction)
	r.PATCH("/transactions/:id", txHandler.UpdateStatus)
	r.DELETE("/transactions/:id", txHandler.DeleteTransaction)
	r.PATCH("/transactions/:id/status", txHandler.UpdateTransactionStatus)

	//graceful shutdown
	srv := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: r,
	}

	go func() {
		log.Println("🚀 Transaction service running at :" + port)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
