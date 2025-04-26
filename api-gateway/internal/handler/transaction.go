package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"api-gateway/internal/client"
)

type TransactionHandler struct {
	BaseURL string
}

func (h TransactionHandler) Create(c *gin.Context) {
	h.proxy(c, http.MethodPost, "/transactions")
}

func (h TransactionHandler) Get(c *gin.Context) {
	h.proxy(c, http.MethodGet, "/transactions/"+c.Param("id"))
}

func (h TransactionHandler) UpdateStatus(c *gin.Context) {
	h.proxy(c, http.MethodPatch, "/transactions/"+c.Param("id")+"/status")
}

func (h TransactionHandler) proxy(c *gin.Context, method, path string) {
	body, _ := io.ReadAll(c.Request.Body)
	resp, err := client.Forward(method, h.BaseURL+path, body, map[string]string{
		"Content-Type":  c.GetHeader("Content-Type"),
		"Authorization": c.GetHeader("Authorization"),
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	err = client.CopyResponse(c.Writer, resp)
	if err != nil {
		return
	}
}

func (h TransactionHandler) Delete(c *gin.Context) {
	h.proxy(c, http.MethodDelete, "/transactions/"+c.Param("id"))
}
