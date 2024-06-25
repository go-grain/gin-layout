package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
	"time"
)

type ErrorRes struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Time    int64  `json:"time"`
}

type Model struct {
	// 自增ID
	ID uint `json:"id" xml:"id" gorm:"primarykey"`
	// 创建时间
	CreatedAt time.Time `json:"createdAt" xml:"createdAt"`
	// 更新时间
	UpdatedAt time.Time `json:"updatedAt" xml:"updatedAt"`
	// 删除时间
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" xml:"deletedAt" gorm:"index"`
}

type MongoModel struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	// 创建时间
	CreatedAt time.Time `json:"createdAt" xml:"createdAt"`
	// 更新时间
	UpdatedAt time.Time `json:"updatedAt" xml:"updatedAt"`
	// 删除时间
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" xml:"deletedAt"`
}

// BeforeCreate 钩子函数： 创建前Gorm会调用
func (m *Model) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return nil
}

// BeforeUpdate 钩子函数： 更新前Gorm会调用
func (m *Model) BeforeUpdate(tx *gorm.DB) (err error) {
	now := time.Now()
	m.UpdatedAt = now
	return nil
}

// PageReq 分页查询使用的结构体
type PageReq struct {
	// 这个字段对前端隐藏,只服务于后端
	UID string ` form:"-" json:"-"`
	// 前端传过来用于查询数据的ID
	ID string ` form:"id" json:"id"`
	// 所查询的数据总量
	Total int64 `json:"total,omitempty"  form:"total"`
	// 页
	Page int `json:"page,omitempty"  form:"page"`
	// 页大小
	PageSize int `json:"pageSize,omitempty"  form:"pageSize"`
	// 查询关键词
	Keyword string `json:"keyword"  form:"keyword"`
	// 查询类型
	Type int `json:"type"  form:"type"`
}
