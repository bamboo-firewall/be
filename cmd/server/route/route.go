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
	router.Use(gin.LoggerWithFormatter(middleware.LogFormatterMiddleware))
	router.GET("/api/v1/ping", handler.Ping)

	{
		hepHandler := handler.NewHEP(service.NewHEP(repo))
		router.POST("/api/v1/hostEndpoints", hepHandler.Create)
		router.GET("/api/v1/hostEndpoints", hepHandler.List)
		router.GET("/api/v1/hostEndpoints/byTenantID/:tenantID/byIP/:ip", hepHandler.Get)
		router.DELETE("/api/v1/hostEndpoints", hepHandler.Delete)

		router.GET("/api/internal/v1/hostEndpoints/fetchPolicies", hepHandler.FetchPolicies)
	}

	{
		gnpHandler := handler.NewGNP(service.NewGNP(repo))
		router.POST("/api/v1/globalNetworkPolicies", gnpHandler.Create)
		router.GET("/api/v1/globalNetworkPolicies", gnpHandler.List)
		router.GET("/api/v1/globalNetworkPolicies/byName/:name", gnpHandler.Get)
		router.DELETE("/api/v1/globalNetworkPolicies", gnpHandler.Delete)
	}

	{
		gnsHandler := handler.NewGNS(service.NewGNS(repo))
		router.POST("/api/v1/globalNetworkSets", gnsHandler.Create)
		router.GET("/api/v1/globalNetworkSets", gnsHandler.List)
		router.GET("/api/v1/globalNetworkSets/byName/:name", gnsHandler.Get)
		router.DELETE("/api/v1/globalNetworkSets", gnsHandler.Delete)
	}

	return router
}
