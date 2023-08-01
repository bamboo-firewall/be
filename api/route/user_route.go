package route

import (
	"time"

	"github.com/bamboo-firewall/be/api/handler"
	"github.com/bamboo-firewall/be/bootstrap"
	"github.com/bamboo-firewall/be/domain"
	"github.com/bamboo-firewall/be/mongo"
	"github.com/bamboo-firewall/be/repository"
	"github.com/bamboo-firewall/be/usecase"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func NewUserRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup, enforcer *casbin.Enforcer) {
	repo := repository.NewUserRepository(db, domain.CollectionUser)
	hl := &handler.UserHandler{
		UserUsecase: usecase.NewUserUsecase(repo, enforcer, timeout),
		Enforcer:    enforcer,
	}
	group.POST("/user/fetch", hl.Fetch)
	group.POST("/user/create", hl.Create)
	group.POST("/user/update", hl.Update)
	group.POST("/user/delete", hl.DeleteById)
}
