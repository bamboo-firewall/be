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

func NewOptionRoute(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup) {
	hepRepo := repository.NewHEPRepository(db, domain.CollectionHEP)
	gnsRepo := repository.NewGNSRepository(db, domain.CollectionGNS)
	policyRepo := repository.NewPolicyRepository(db, domain.CollectionPolicy)
	hl := &handler.OptionHandler{
		HEPUsecase:    usecase.NewHEPUsecase(hepRepo, timeout),
		GNSUsecase:    usecase.NewGNSUsecase(gnsRepo, timeout),
		PolicyUsecase: usecase.NewPolicyUsecase(policyRepo, timeout),
	}
	group.POST("/options/fetch", hl.FetchByType)
}
