package route

import (
	"time"

	"github.com/bamboo-firewall/be/api/handler"
	"github.com/bamboo-firewall/be/bootstrap"
	"github.com/bamboo-firewall/be/domain"
	"github.com/bamboo-firewall/be/mongo"
	"github.com/bamboo-firewall/be/repository"
	"github.com/bamboo-firewall/be/usecase"
	"github.com/gin-gonic/gin"
)

func NewGNSRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	gr := repository.NewGNSRepository(db, domain.CollectionGNS)
	gc := &handler.GNSHandler{
		GNSUsecase: usecase.NewGNSUsecase(gr, timeout),
	}
	group.POST("/gns/fetch", gc.Fetch)
	group.POST("/gns/search", gc.Search)
}
