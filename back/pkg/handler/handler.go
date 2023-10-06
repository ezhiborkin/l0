package handler

import (
	"back/pkg/cashe"
	"back/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
	cashe    *cashe.Cache
}

func NewHandler(services *service.Service, cashe *cashe.Cache) *Handler {
	return &Handler{
		services: services,
		cashe:    cashe,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		cashe := api.Group("/cache")
		{
			cashe.GET("/:id", h.getByIdCache)
		}
		bd := api.Group("/bd")
		{
			bd.GET("/:id", h.getById)
		}

	}

	lol := router.Group("/lol")
	lol.GET("/lol", func(c *gin.Context) {
		c.JSON(http.StatusOK, 123)
	})

	return router
}
