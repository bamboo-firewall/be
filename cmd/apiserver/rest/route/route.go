package route

import (
	"net/http"

	"github.com/bamboo-firewall/be/cmd/apiserver/rest/middleware"
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
