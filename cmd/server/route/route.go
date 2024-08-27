package route

import (
	"github.com/bamboo-firewall/be/api/v1/handler"
	"github.com/bamboo-firewall/be/cmd/server/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandler() http.Handler {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// define api
	{

	}
	router.Handle(http.MethodGet, "/api/v1/ping", handler.Ping)

	return router
}
