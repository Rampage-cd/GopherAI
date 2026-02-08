package mysql

import (
	"GopherAI/config"
	"GopherAI/model"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

//全局的数据库连接对象（数据库连接池的封装）

func InitMysql() error {
	host := config.GetConfig().MysqlHost
	port := config.GetConfig().MysqlPort
	dbname := config.GetConfig().MysqlDatabaseName
	username := config.GetConfig().MysqlUser
	password := config.GetConfig().MysqlPassword
	charset := config.GetConfig().MysqlCharset

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local", username, password, host, port, dbname, charset)
	//拼接MySQL的DSN（Data Source Name）
	//用户名:密码@tcp(主机:端口)/数据库名?参数

	var log logger.Interface
	if gin.Mode() == "debug" {
		log = logger.Default.LogMode(logger.Info)
	} else {
		log = logger.Default
	} //debug模式：打印SQL语句
	//release模式：减少日志输出

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: log,
	}) //初始化数据库连接池(使用gorm)
	if err != nil {
		return err
	}

	sqlDB, err := db.DB() //获取底层（原生连接池）对象
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(10)           //最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          //最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) //单个连接最大存活时间

	DB = db
	//保存到全局连接对象

	return migration()
	//自动建表
}

func migration() error {
	return DB.AutoMigrate(
		new(model.User),
		new(model.Session),
		new(model.Message),
	) //如果表不存在，则创建表(三个表，用户表，会话表，信息表)
	//如果字段不存在，则添加字段
}

func InsertUser(user *model.User) (*model.User, error) {
	err := DB.Create(&user).Error
	//获取错误信息
	return user, err
} //插入用户

func GetUserByUsername(username string) (*model.User, error) {
	user := new(model.User)
	err := DB.Where("username = ?", username).First(user).Error
	//相当于sql中的SELECT * FROM users WHERE username = ? LIMIT 1;
	//?是为了防止SQL注入
	return user, err
} //根据用户名查询用户
