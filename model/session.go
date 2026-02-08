package model

import (
	"time"

	"gorm.io/gorm"
)

// 数据库模型
type Session struct {
	ID        string         `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserName  string         `gorm:"index;not null" json:"username"`
	Title     string         `gorm:"type:varchar(100)" json:"title"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	//软删除，给deleted_at字段赋值时间，正常查询时自动过滤掉
	//json:"-"表示JSON序列化时忽略该字段
}

// 接口返回模型
type SessionInfo struct {
	SessionID string `json:"sessionId"`
	Title     string `json:"name"`
} //不是数据库表模型，更贴近前端使用
