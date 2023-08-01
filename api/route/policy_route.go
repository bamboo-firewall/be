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

func NewPolicyRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	r := repository.NewPolicyRepository(db, domain.CollectionPolicy)
	h := &handler.PolicyHandler{
		PolicyUsecase: usecase.NewPolicyUsecase(r, timeout),
	}
	group.POST("/policy/fetch", h.Fetch)
	group.POST("/policy/search", h.Search)
}
