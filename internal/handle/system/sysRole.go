package handler

import (
	"github.com/gin-gonic/gin"
	service "github.com/go-grain/gin-layout/internal/service/system"
	model "github.com/go-grain/gin-layout/model/system"
	"github.com/go-grain/gin-layout/pkg/convert"
	"github.com/go-grain/gin-layout/pkg/response"
)

type RoleHandle struct {
	res response.Response
	sv  *service.RoleService
}

func NewRoleHandle(sv *service.RoleService) *RoleHandle {
	return &RoleHandle{
		sv: sv,
	}
}

func (r *RoleHandle) InitRole() error {
	return r.sv.InitRole()
}

func (r *RoleHandle) CreateRole(ctx *gin.Context) {
	res := r.res.New()
	role := model.CreateSysRole{}
	err := ctx.ShouldBindJSON(&role)
	if err != nil {
		res.WithCode(4000).WithMessage(err.Error()).Fail(ctx)
		return
	}
	err = r.sv.CreateRole(&role, ctx)
	if err != nil {
		res.WithCode(4000).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("创建角色成功").Success(ctx)
}

func (r *RoleHandle) GetRoleList(ctx *gin.Context) {
	res := r.res.New()
	req := model.SysRoleQueryPage{}
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		res.WithCode(4000).WithMessage("参数解析失败").Fail(ctx)
		return
	}
	list, err := r.sv.GetRoleList(&req, ctx)
	if err != nil {
		res.WithCode(4000).WithMessage(err.Error()).Fail(ctx)
		return
	}
	res.WithMessage("成功").WithData(list).Success(ctx)
}

func (r *RoleHandle) UpdateRole(ctx *gin.Context) {
	res := r.res.New()
	role := model.SysRole{}
	err := ctx.ShouldBindJSON(&role)
	if err != nil {
		res.WithCode(4000).WithMessage("解析参数失败").Fail(ctx)
		return
	}
	err = r.sv.UpdateRole(&role, ctx)
	if err != nil {
		res.WithCode(4000).WithMessage("更新角色失败").Fail(ctx)
		return
	}
	res.WithMessage("更新角色成功").Success(ctx)
}

func (r *RoleHandle) DeleteRoleById(ctx *gin.Context) {
	reply := r.res.New()
	role := convert.String2Int(ctx.Query("id"))
	if role == 0 {
		reply.WithCode(4000).WithMessage("ID不能为空").Fail(ctx)
		return
	}
	err := r.sv.DeleteRoleById(uint(role), ctx)
	if err != nil {
		reply.WithCode(4000).WithMessage("删除角色失败").Fail(ctx)
		return
	}
	reply.WithMessage("删除角色成功").Success(ctx)
}

func (r *RoleHandle) DeleteRoleByIds(ctx *gin.Context) {
	reply := r.res.New()
	api := struct {
		Roles []uint `json:"ids"`
	}{}
	err := ctx.ShouldBindJSON(&api)
	if err != nil {
		reply.WithCode(4000).WithMessage("ids不能为空").Fail(ctx)
		return
	}
	err = r.sv.DeleteRoleByIds(api.Roles, ctx)
	if err != nil {
		reply.WithCode(4000).WithMessage("批量删除角色失败").Fail(ctx)
		return
	}
	reply.WithMessage("批量删除角色成功").Success(ctx)
}
