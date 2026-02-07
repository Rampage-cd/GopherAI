package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

// 当 toml 库解析配置文件时：
//
//	port = 8080
//
// 会被映射到 Port 这个字段
//
// 如果不写这个 tag，toml 默认会尝试找名为 "Port" 的配置项
type MainConfig struct {
	Port    int    `toml:"port"`
	AppName string `toml:"appName"`
	Host    string `toml:"host"`
}

//`toml:"port"`等类似的内容成为结构体标签，Go编译器不会解释他的含义，他的含义由使用反射的库（toml库）来定义
//toml是一种轻量级的，高度可读的且语义明确的配置文件格式（大小写敏感）
//之所以要写结构体标签，是因为方便toml配置文件查找
//结构体标签的统一格式为`toml:"port" json:"port" yaml:"port"`

type EmailConfig struct {
	Authcode string `toml:"authcode"`
	Email    string `toml:"email"`
} //邮件配置

type RedisConfig struct {
	RedisPort     int    `toml:"port"`
	RedisDb       int    `toml:"db"`
	RedisHost     string `toml:"host"`
	RedisPassword string `toml:"password"`
}

type MysqlConfig struct {
	MysqlPort         int    `toml:"port"`
	MysqlHost         string `toml:"host"`
	MysqlUser         string `toml:"user"`
	MysqlPassword     string `toml:"password"`
	MysqlDatabaseName string `toml:"databaseName"`
	MysqlCharset      string `toml:"charset"`
}

type JwtConfig struct {
	ExpireDuration int    `toml:"expire_duration"`
	Issuer         string `toml:"issuer"`
	Subject        string `toml:"subject"`
	Key            string `toml:"key"`
} //JSON Web Token配置

type Rabbitmq struct {
	RabbitmqPort     int    `toml:"port"`
	RabbitmqHost     string `toml:"host"`
	RabbitmqUsername string `toml:"username"`
	RabbitmqPassword string `toml:"password"`
	RabbitmqVhost    string `toml:"vhost"`
} //消息队列配置

type Config struct {
	EmailConfig `toml:"emailConfig"`
	RedisConfig `toml:"redisConfig"`
	MysqlConfig `toml:"mysqlConfig"`
	JwtConfig   `toml:"jwtConfig"`
	MainConfig  `toml:"mainConfig"`
	Rabbitmq    `toml:"rabbitmqConfig"`
} //结构体嵌套，子结构体Config可以直接使用父结构体的字段和方法

type RedisKeyConfig struct {
	CaptchaPrefix string
}

var config *Config

// InitConfig 初始化项目配置
func InitConfig() error {
	// 设置配置文件路径（相对于 main.go 所在的目录）
	if _, err := toml.DecodeFile("config/config.toml", config); err != nil {
		log.Fatal(err.Error())
		return err
	} //这个操作相当于给结构体中的所有字段赋值
	return nil
	//1. toml 读取 config.toml
	//2. 通过反射遍历 Config 的字段
	//3. 对每个字段：
	//- 读取 toml tag
	//- 用 tag 的值作为 key 去配置文件中找
	//4. 找到后，把值写入 struct 对应字段
}

func GetConfig() *Config {
	if config == nil {
		config = new(Config)
		_ = InitConfig()
	}
	return config
} //获取配置好的文件（先创建，后初始化）
