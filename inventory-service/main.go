package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/shoe-store/inventory-service/internal/natsadapter"

	"github.com/shoe-store/inventory-service/internal/cache"
	"github.com/shoe-store/inventory-service/internal/config"
	"github.com/shoe-store/inventory-service/internal/repository"
	"github.com/shoe-store/inventory-service/internal/server"
	"github.com/shoe-store/inventory-service/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var requestCount = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "shoestore_requests_total",
		Help: "Total number of HTTP requests received",
	},
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	go func() {
		// Prometheus metrics
		r := gin.Default()

		prometheus.MustRegister(requestCount)

		r.Use(func(c *gin.Context) {
			requestCount.Inc()
			c.Next()
		})

		r.GET("/metrics", gin.WrapH(promhttp.Handler()))

		err := r.Run(":5053")
		if err != nil {
			return
		}
	}()
	// Initialize Redis
	cache.InitRedis()

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Printf("Error disconnecting MongoDB: %v", err)
		} else {
			log.Println("MongoDB disconnected")
		}
	}(client, context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Successfully connected to MongoDB")

	// Initialize repository
	repo := repository.NewMongoRepository(client.Database("OnlineStore"))

	// Initialize service
	svc := service.NewInventoryService(repo)

	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	natsadapter.SubscribeToOrderCreated(nc, svc)

	// Initialize server
	srv := server.NewServer(svc)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting gRPC server on port %s", cfg.Port)
		if err := srv.Start(cfg.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	log.Println("Shutting down server...")
	srv.Stop()

	if err := client.Disconnect(context.Background()); err != nil {
		log.Printf("Error disconnecting MongoDB: %v", err)
	}
}
