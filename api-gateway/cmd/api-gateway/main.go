package main

import (
	client "api-gateway/internal/client/grpc"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"api-gateway/internal/config"
	"api-gateway/internal/handler"
)

var requestCount = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "shoestore_requests_total",
		Help: "Total number of HTTP requests received",
	},
)

func main() {
	cfg := config.Load()

	client.InitTransactionGRPCClient()
	client.InitAuthGRPCClient()
	client.InitInventoryGRPCClient()

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Prometheus metrics
	prometheus.MustRegister(requestCount)
	r.Use(func(c *gin.Context) {
		requestCount.Inc()
		c.Next()
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Protected group
	protected := r.Group("/")

	//users
	auth := handler.AuthHandler{}
	r.POST("/users/signup", auth.SignUp)
	r.GET("/users/confirm", auth.ConfirmEmail)
	r.POST("/users/login", auth.Login)
	r.POST("/users/change-password", auth.ChangePassword)

	// Transactions
	tx := handler.TransactionHandler{}
	protected.POST("/transactions", tx.Create)
	protected.GET("/transactions/:id", tx.Get)
	protected.PATCH("/transactions/:id/status", tx.UpdateStatus)
	protected.DELETE("/transactions/:id", tx.Delete)

	// Inventory routes
	inv := handler.InventoryHandler{}
	protected.GET("/inventory/sneakers", inv.GetSneakers)
	protected.POST("/inventory/sneakers", inv.CreateSneaker)
	protected.PUT("/inventory/sneakers", inv.EditSneaker)
	protected.DELETE("/inventory/sneakers", inv.RemoveSneaker)
	protected.GET("/inventory/public-sneakers", inv.GetPublicSneakers)

	srv := &http.Server{
		Addr:           "0.0.0.0:" + cfg.Port,
		Handler:        r,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("API‑Gateway listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("gateway stopped: %v", err)
	}
}
