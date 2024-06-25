package core

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-grain/gin-layout/config"
	"github.com/go-grain/gin-layout/internal/repo/data"
	"github.com/go-grain/gin-layout/internal/repo/system/query"
	router "github.com/go-grain/gin-layout/internal/router/system"
	"github.com/go-grain/gin-layout/log"
	"github.com/go-grain/gin-layout/middleware"
	redisx "github.com/go-grain/gin-layout/pkg/redis"
	"github.com/go-grain/gin-layout/pkg/response"
	"gorm.io/gorm"
	"os"
)

var (
	Name    string
	Version string
	id, _   = os.Hostname()
)

type IInit interface {
	init(grain *Great) error
}

type Great struct {
	db       *gorm.DB
	sysLog   log.Logger
	engine   *gin.Engine
	conf     *config.Config
	rdb      redisx.IRedis
	enforcer *casbin.CachedEnforcer
}

type InitConf struct{}

func (InitConf) init(great *Great) (err error) {
	great.conf, err = config.InitConfig()
	if err != nil {
		return
	}

	os.Mkdir(".tmp/", 0o664)
	file, err := os.OpenFile(".tmp/great.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o664)
	if err != nil {
		return err
	}

	great.sysLog = log.With(log.NewStdLogger(file),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
	)

	great.db, err = data.InitDB(*great.conf)
	if err != nil {
		return
	}

	great.rdb, err = data.InitRedis(*great.conf)
	if err != nil {
		return
	}
	return
}

type InitRouter struct{}

func (InitRouter) init(great *Great) (err error) {
	great.engine = gin.Default()
	gin.SetMode(great.conf.Gin.Model)
	great.engine.Use(middleware.Cors())

	routerGroup := great.engine.Group("api/v1")
	great.engine.NoRoute(func(ctx *gin.Context) {
		reply := response.Response{}
		reply.WithCode(404).WithMessage("请求路径不正确").Fail(ctx)
	})

	router.NewRoleRouter(routerGroup, great.rdb, great.conf, great.sysLog).InitRouters().InitRole()
	router.NewSysUserRouter(great.engine, routerGroup, great.rdb, great.conf, great.sysLog).InitRouters().InitUser()

	return nil
}

type RunGin struct{}

func (RunGin) init(great *Great) (err error) {
	if err = great.engine.Run(great.conf.Gin.Host); err != nil {
		return err
	}
	return nil
}

type InitGenQuery struct{}

func (InitGenQuery) init(great *Great) (err error) {
	query.SetDefault(great.db)
	return nil
}

func (Great) Do(great *Great, init []IInit) {
	for _, iInit := range init {
		err := iInit.init(great)
		if err != nil {
			panic(err)
		}
	}
}

func Init() {
	great := Great{}
	init := []IInit{
		&InitConf{},
		&InitGenQuery{},
		&InitRouter{},
		&RunGin{},
	}

	great.Do(&great, init)
}
