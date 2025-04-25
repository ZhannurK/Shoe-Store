package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"transaction-service/internal/handlers"
	"transaction-service/internal/repositories"
	"transaction-service/internal/usecase"
)

func main() {
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

	mongoClient := repositories.InitMongo(mongoURI)
	txRepo := repositories.NewTransactionRepository(mongoClient, "OnlineStore")
	txUC := usecase.NewTransactionUseCase(txRepo)
	txHandler := handlers.NewTransactionHandler(txUC)

	r := gin.Default()

	r.POST("/transactions", txHandler.CreateTransaction)
	r.GET("/transactions/:id", txHandler.GetTransaction)
	r.PATCH("/transactions/:id", txHandler.UpdateStatus)
	r.DELETE("/transactions/:id", txHandler.DeleteTransaction)
	r.PATCH("/transactions/:id/status", txHandler.UpdateTransactionStatus)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Println("🚀 Transaction service running at :" + port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
