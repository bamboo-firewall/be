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

func NewSignupRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup, enforcer *casbin.Enforcer) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)
	sc := handler.SignupHandler{
		SignupUsecase: usecase.NewSignupUsecase(ur, timeout),
		Enforcer:      enforcer,
		Env:           env,
	}
	group.POST("/signup", sc.Signup)
}
