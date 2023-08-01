package route

import (
	"net/http"
	"time"

	"github.com/bamboo-firewall/be/bootstrap"
	"github.com/bamboo-firewall/be/mongo"
	"github.com/gin-gonic/gin"
)

func NewPingRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	group.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})
}
