package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getById(c *gin.Context) {
	orderID := c.Param("id")
	fmt.Println(orderID)
	order, err := h.services.GetById(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to get order by ID"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *Handler) getByIdCache(c *gin.Context) {
	orderID := c.Param("id")
	fmt.Println(orderID)
	order, err := h.cashe.Get(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed to get order by ID"})
		return
	}

	c.JSON(http.StatusOK, order)
}
