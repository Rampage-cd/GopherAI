package model

import (
	"time"
)

type Message struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SessionID string    `gorm:"index;not null;type:varchar(36)" json:"session_id"`
	UserName  string    `gorm:"type:varchar(20)" json:"username"` //type为数据库列类型
	Content   string    `gorm:"type:text" json:"content"`
	IsUser    bool      `gorm:"not null;" json:"is_user"`
	CreatedAt time.Time `json:"created_at"`
	//struct tag不会被编译器看见，而是由“使用这个struct的库”去解析
}

// json用于在序列化和反序列化时，告知字段名用什么
// gorm用于定义数据库映射规则
type History struct {
	IsUser  bool   `json:"is_user"`
	Content string `json:"content"`
}
