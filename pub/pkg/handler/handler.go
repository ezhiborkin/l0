package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/stan.go"
)

type Handler struct {
	sc stan.Conn
}

func NewHandler(sc stan.Conn) *Handler {
	return &Handler{sc: sc}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(cors.Default())

	api := router.Group("/api")
	{
		api.POST("/publish", h.publish)
		// api.POST("/kek", h.kek)
	}

	return router

}
