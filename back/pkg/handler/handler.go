package handler

import (
	"back/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		orders := api.Group("/orders")
		{
			orders.GET("/:id", func(c *gin.Context) {
				order, err := h.getById(c)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order"})
					return
				}
				c.JSON(http.StatusOK, order)
			})
		}
	}

	router.GET("/lol", func(c *gin.Context) {
		c.JSON(http.StatusOK, 123)
	})

	return router
}
