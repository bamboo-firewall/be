package route

import (
	"fmt"
	"time"

	"github.com/bamboo-firewall/be/api/middleware"
	"github.com/bamboo-firewall/be/bootstrap"
	"github.com/bamboo-firewall/be/mongo"
	"github.com/casbin/casbin/v2"

	mongodbadapter "github.com/casbin/mongodb-adapter/v3"

	"github.com/gin-gonic/gin"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db mongo.Database, gin *gin.Engine) {
	gin.Use(middleware.CORSMiddleware(env))

	adapterConfig := mongodbadapter.AdapterConfig{
		DatabaseName: env.DBName,
	}

	adapter, err := mongodbadapter.NewAdapterByDB(db.Client().MongoClient(), &adapterConfig)

	if err != nil {
		panic(err)
	}

	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		panic(fmt.Sprintf("failed to create casbin enforcer: %v", err))
	}

	//add policy
	if hasPolicy := enforcer.HasPolicy("admin", "user", "write"); !hasPolicy {
		enforcer.AddPolicy("admin", "user", "write")
	}
	if hasPolicy := enforcer.HasPolicy("admin", "user", "read"); !hasPolicy {
		enforcer.AddPolicy("admin", "user", "read")
	}
	if hasPolicy := enforcer.HasPolicy("devops", "user", "read"); !hasPolicy {
		enforcer.AddPolicy("devops", "user", "read")
	}

	publicRouter := gin.Group("api/")
	// All Public APIs
	NewPingRouter(env, timeout, db, publicRouter)
	NewSignupRouter(env, timeout, db, publicRouter, enforcer)
	NewLoginRouter(env, timeout, db, publicRouter)
	NewRefreshTokenRouter(env, timeout, db, publicRouter)

	// Privated APIs
	protectedRouter := gin.Group("api/v1/")
	protectedRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	NewProfileRouter(env, timeout, db, protectedRouter, enforcer)
	NewOptionRoute(env, timeout, db, protectedRouter)
	NewHEPRouter(env, timeout, db, protectedRouter)
	NewGNSRouter(env, timeout, db, protectedRouter)
	NewPolicyRouter(env, timeout, db, protectedRouter)
	NewStatisticRouter(env, timeout, db, protectedRouter)

	// Admin APIs
	adminRouter := gin.Group("api/v1/admin/")
	adminRouter.Use(middleware.JwtAuthMiddleware(env.AccessTokenSecret))
	adminRouter.Use(middleware.Authorize("user", "write", enforcer))
	NewUserRouter(env, timeout, db, adminRouter, enforcer)
}
