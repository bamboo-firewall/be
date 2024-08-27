package route

import (
	"github.com/bamboo-firewall/be/v1/middleware"
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

	return router
}
