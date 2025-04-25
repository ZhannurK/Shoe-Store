package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"transaction-service/internal/domain"
	"transaction-service/internal/usecase"
)

type TransactionHandler struct {
	UC *usecase.TransactionUseCase
}

func NewTransactionHandler(uc *usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{UC: uc}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var tx domain.Transaction
	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx.TransactionID = primitive.NewObjectID().Hex()
	tx.Status = domain.StatusPending
	tx.CreatedAt = time.Now()
	tx.UpdatedAt = time.Now()

	if err := h.UC.Create(c.Request.Context(), &tx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	c.JSON(http.StatusCreated, tx)
}

func (h *TransactionHandler) UpdateStatus(c *gin.Context) {
	txID := c.Param("id")
	var body struct {
		Status domain.TransactionStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.UC.UpdateStatus(c.Request.Context(), txID, body.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction status updated"})
}

func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	txID := c.Param("id")
	tx, err := h.UC.GetByID(c.Request.Context(), txID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	c.JSON(http.StatusOK, tx)
}

func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	id := c.Param("id")
	err := h.UC.DeleteTransaction(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted"})
}

type statusPayload struct {
	Status string `json:"status" binding:"required"`
}

func (h *TransactionHandler) UpdateTransactionStatus(c *gin.Context) {
	id := c.Param("id")

	var p statusPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newStatus domain.TransactionStatus
	switch p.Status {
	case string(domain.StatusPaid):
		newStatus = domain.StatusPaid
	case string(domain.StatusDeclined):
		newStatus = domain.StatusDeclined
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	if err := h.UC.UpdateStatus(c, id, newStatus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "status updated", "status": newStatus})
}
