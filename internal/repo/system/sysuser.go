package repo

import (
	"github.com/go-grain/gin-layout/internal/repo/system/query"
	service "github.com/go-grain/gin-layout/internal/service/system"
	model "github.com/go-grain/gin-layout/model/system"
	redisx "github.com/go-grain/gin-layout/pkg/redis"
)

type SysUserRepo struct {
	rdb   redisx.IRedis
	query *query.Query
}

func NewSysUserRepo(rdb redisx.IRedis) service.ISysUserRepo {
	return &SysUserRepo{
		rdb:   rdb,
		query: query.Q,
	}
}

func (r *SysUserRepo) Login(user *model.LoginReq) (*model.SysUser, error) {
	userinfo, err := r.query.SysUser.Where(r.query.SysUser.Username.Eq(user.Username)).First()
	if err != nil {
		return nil, err
	}
	return userinfo, err
}
