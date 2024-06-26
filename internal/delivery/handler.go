package delivery

import (
	"github.com/gin-gonic/gin"
	"psPro-task/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{services: s}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	api := router.Group("/api")
	{
		api.POST("", h.createCommand)
		api.GET("", h.listCommands)
		api.GET("/comm", h.oneCommand)
		api.PATCH("/stop/:id", h.stopCommand)
		api.PATCH("/start/:id", h.startCommand)
		api.PATCH("/kill/:id", h.killCommand)
	}
	return router
}
