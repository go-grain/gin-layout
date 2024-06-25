package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-grain/gin-layout/config"
	handler "github.com/go-grain/gin-layout/internal/handle/system"
	repo "github.com/go-grain/gin-layout/internal/repo/system"
	service "github.com/go-grain/gin-layout/internal/service/system"
	"github.com/go-grain/gin-layout/log"
	"github.com/go-grain/gin-layout/middleware"
	redisx "github.com/go-grain/gin-layout/pkg/redis"
)

type RoleRouter struct {
	api     *handler.RoleHandle
	public  gin.IRoutes
	private gin.IRoutes
}

func NewRoleRouter(routerGroup *gin.RouterGroup, rdb redisx.IRedis, conf *config.Config, logger log.Logger) *RoleRouter {
	data := repo.NewRoleRepo(rdb)
	sv := service.NewRoleService(data, rdb, conf, logger)
	return &RoleRouter{
		api:    handler.NewRoleHandle(sv),
		public: routerGroup.Group("sysRole"),
		private: routerGroup.Group("sysRole").Use(
			middleware.JwtAuth(*conf, rdb),
		),
	}
}

func (r *RoleRouter) InitRouters() *RoleRouter {
	r.private.POST("", r.api.CreateRole)
	r.private.PUT("", r.api.UpdateRole)
	r.private.GET("list", r.api.GetRoleList)
	r.private.DELETE("", r.api.DeleteRoleById)
	r.private.DELETE("deleteRoleByIds", r.api.DeleteRoleByIds)
	return r
}

func (r *RoleRouter) InitRole() *RoleRouter {
	_ = r.api.InitRole()
	return r
}
