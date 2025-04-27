package handler

import (
	grpcClient "api-gateway/internal/client/grpc"
	pb "api-gateway/proto/transaction"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TransactionHandler struct{}

func (h TransactionHandler) Create(c *gin.Context) {
	var body struct {
		UserID    string        `json:"userId"`
		CartItems []pb.CartItem `json:"cartItems"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx, err := grpcClient.GRPCCreateTransaction(body.UserID, pbToPtr(body.CartItems))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, tx)
}

func (h TransactionHandler) Get(c *gin.Context) {
	tx, err := grpcClient.GRPCGetTransaction(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tx)
}

func (h TransactionHandler) UpdateStatus(c *gin.Context) {
	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := grpcClient.GRPCUpdateStatus(c.Param("id"), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h TransactionHandler) Delete(c *gin.Context) {
	if err := grpcClient.GRPCDeleteTransaction(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func pbToPtr(in []pb.CartItem) []*pb.CartItem {
	out := make([]*pb.CartItem, len(in))
	for i := range in {
		out[i] = &in[i]
	}
	return out
}

//package handler
//
//import (
//    "net/http"
//    "github.com/gin-gonic/gin"
//    grpcClient "api-gateway/internal/client/grpc"
//    pb "api-gateway/proto/transaction"
//)
//
//type TransactionHandler struct{}
//
//func (h TransactionHandler) Create(c *gin.Context) {
//    var body struct {
//        UserID    string        `json:"userId"`
//        CartItems []pb.CartItem `json:"cartItems"`
//    }
//    if err := c.ShouldBindJSON(&body); err != nil {
//        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//        return
//    }
//
//    tx, err := grpcClient.GRPCCreateTransaction(body.UserID, pbToPtr(body.CartItems))
//    if err != nil {
//        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//        return
//    }
//    c.JSON(http.StatusCreated, tx)
//}
//
//func (h TransactionHandler) Get(c *gin.Context) {
//    tx, err := grpcClient.GRPCGetTransaction(c.Param("id"))
//    if err != nil {
//        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//        return
//    }
//    c.JSON(http.StatusOK, tx)
//}
//
//func (h TransactionHandler) UpdateStatus(c *gin.Context) {
//    var req struct{ Status string `json:"status"` }
//    if err := c.ShouldBindJSON(&req); err != nil {
//        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//        return
//    }
//    if err := grpcClient.GRPCUpdateStatus(c.Param("id"), req.Status); err != nil {
//        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//        return
//    }
//    c.Status(http.StatusOK)
//}
//
//func (h TransactionHandler) Delete(c *gin.Context) {
//    if err := grpcClient.GRPCDeleteTransaction(c.Param("id")); err != nil {
//        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//        return
//    }
//    c.Status(http.StatusNoContent)
//}
//
//// helper – converts []T to []*T without extra alloc
//func pbToPtr(in []pb.CartItem) []*pb.CartItem {
//    out := make([]*pb.CartItem, len(in))
//    for i := range in {
//        out[i] = &in[i]
//    }
//    return out
//}
