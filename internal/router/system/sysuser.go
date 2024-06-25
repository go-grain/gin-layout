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

type SysUserRouter struct {
	rdb             redisx.IRedis
	api             *handler.SysUserHandle
	engine          *gin.Engine
	public          gin.IRoutes
	private         gin.IRoutes
	privateRoleAuth gin.IRoutes
}

func NewSysUserRouter(engine *gin.Engine, routerGroup *gin.RouterGroup, rdb redisx.IRedis, conf *config.Config, logger log.Logger) *SysUserRouter {
	data := repo.NewSysUserRepo(rdb)
	sv := service.NewSysUserService(data, rdb, conf, logger)
	return &SysUserRouter{
		rdb:    rdb,
		api:    handler.NewSysUserHandle(sv),
		engine: engine,
		public: routerGroup.Group("sysUser"),
		private: routerGroup.Group("sysUser").Use(
			middleware.JwtAuth(*conf, rdb)),
		privateRoleAuth: routerGroup.Group("sysUser").Use(
			middleware.JwtAuth(*conf, rdb),
		),
	}
}

func (r *SysUserRouter) InitRouters() *SysUserRouter {
	//登录接口
	r.public.POST("login", r.api.Login)
	return r
}

func (r *SysUserRouter) InitUser() *SysUserRouter {
	_ = r.api.InitUser()
	return r
}
