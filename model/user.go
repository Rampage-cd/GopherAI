package model

//model包的作用是用go的结构体，描述数据库中的“表结构”
//gorm会根据这些结构体，以及后面的tag自动建表和自动生成SQL，并自动把查询结果映射回来

//user表示系统的使用者
//session表示用户的一次会话
//message表示会话中的一条具体信息
import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(50)" json:"name"`
	Email     string         `gorm:"type:varchar(100);index" json:"email"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex" json:"username"` // 唯一索引
	Password  string         `gorm:"type:varchar(255)" json:"-"`                   // 不返回给前端
	CreatedAt time.Time      `json:"created_at"`                                   // 自动时间戳
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 支持软删除
}
