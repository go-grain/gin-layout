package model

// SysRole 角色结构体
type SysRole struct {
	Model
	Role     string `form:"role" json:"role" xml:"role" gorm:"unique;not null;comment:角色ID" binding:"required"`
	RoleName string `form:"roleName" json:"roleName" xml:"roleName" gorm:"unique;not null;comment:角色名称" binding:"required"`
}

func (SysRole) TableName() string {
	return "sys_roles"
}

type CreateSysRole struct {
	Role     string `json:"role" binding:"required"`
	RoleName string `json:"roleName" binding:"required"`
}

type SysRoleQueryPage struct {
	PageReq
	Role     string `json:"role" form:"role"`
	RoleName string `json:"RoleName" form:"RoleName"`
}
