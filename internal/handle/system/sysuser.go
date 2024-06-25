package handler

import (
	"github.com/gin-gonic/gin"
	service "github.com/go-grain/gin-layout/internal/service/system"
	model "github.com/go-grain/gin-layout/model/system"
	"github.com/go-grain/gin-layout/pkg/response"
)

type SysUserHandle struct {
	res response.Response
	sv  *service.SysUserService
}

func NewSysUserHandle(sv *service.SysUserService) *SysUserHandle {
	return &SysUserHandle{
		sv: sv,
	}
}

func (r *SysUserHandle) InitUser() error {
	return r.sv.InitSysUser()
}

func (r *SysUserHandle) Login(ctx *gin.Context) {
	reply := r.res.New()
	user := model.LoginReq{}
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		reply.WithCode(4000).WithMessage("请求数据有误").Fail(ctx)
		return
	}
	token, err := r.sv.Login(&user, ctx)
	if err != nil {
		reply.WithCode(4000).WithMessage(err.Error()).Fail(ctx)
		return
	}
	reply.WithMessage("欢迎回来").WithData(gin.H{"token": token}).Success(ctx)
}
