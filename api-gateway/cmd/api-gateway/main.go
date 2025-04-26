package main

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"time"

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

	// Transactions
	tx := handler.TransactionHandler{BaseURL: cfg.TransactionURL}
	protected.POST("/transactions", tx.Create)
	protected.GET("/transactions/:id", tx.Get)
	protected.PATCH("/transactions/:id/status", tx.UpdateStatus)
	protected.DELETE("/transactions/:id", tx.Delete)

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
