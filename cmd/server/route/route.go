package route

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bamboo-firewall/be/api/v1/handler"
	"github.com/bamboo-firewall/be/cmd/server/middleware"
	"github.com/bamboo-firewall/be/cmd/server/pkg/repository"
	"github.com/bamboo-firewall/be/domain/service"
)

func RegisterHandler(repo *repository.PolicyDB) http.Handler {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.GET("/api/v1/ping", handler.Ping)

	{
		hepHandler := handler.NewHEP(service.NewHEP(repo))
		router.POST("/api/v1/hostEndpoints", hepHandler.Create)
		router.DELETE("/api/v1/hostEndpoints", hepHandler.Delete)

		router.POST("/api/internal/v1/hostEndpoints/byName/:name/fetchPolicies", hepHandler.FetchPolicies)
	}

	{
		gnpHandler := handler.NewGNP(service.NewGNP(repo))
		router.POST("/api/v1/globalNetworkPolicies", gnpHandler.Create)
		router.DELETE("/api/v1/globalNetworkPolicies", gnpHandler.Delete)
	}

	{
		gnsHandler := handler.NewGNS(service.NewGNS(repo))
		router.POST("/api/v1/globalNetworkSets", gnsHandler.Create)
		router.DELETE("/api/v1/globalNetworkSets", gnsHandler.Delete)
	}

	return router
}
