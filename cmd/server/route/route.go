package route

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bamboo-firewall/be/api/v1/handler"
	"github.com/bamboo-firewall/be/cmd/server/middleware"
	"github.com/bamboo-firewall/be/cmd/server/pkg/repository"
	"github.com/bamboo-firewall/be/domain/service"
)

func RegisterHandler(repo *repository.PolicyMongo) http.Handler {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.GET("/api/v1/ping", handler.Ping)

	{
		hepHandler := handler.NewHEP(service.NewHEP(repo))
		router.POST("/api/v1/hostEndpoints", hepHandler.Create)
		router.DELETE("/api/v1/hostEndpoints", hepHandler.Delete)
	}

	return router
}
