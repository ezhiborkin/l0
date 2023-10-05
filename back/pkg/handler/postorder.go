package handler

import (
	"back/pkg/order"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getById(c *gin.Context) (*order.OrderData, error) {
	orderId := c.Param("id")

	order, err := h.services.GetById(orderId)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return nil, err
	}

	return order, nil
}
