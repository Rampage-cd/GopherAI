package redis

import (
	"GopherAI/config"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

// Rdb 是一个全局Redis客户端实例
// 整个项目中复用这一份连接池
var Rdb *redis.Client

var ctx = context.Background()

func Init() {
	conf := config.GetConfig()
	host := conf.RedisConfig.RedisHost
	port := conf.RedisConfig.RedisPort
	password := conf.RedisConfig.RedisPassword
	db := conf.RedisDb
	addr := host + ":" + strconv.Itoa(port)
	//组装Redis地址，例如"127.0.0.1:6379"

	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,     //Redis地址
		Password: password, //Redis密码
		DB:       db,       //Redis数据库编号
	})
	//创建Redis客户端

}

// 为某个邮箱设置验证码
func SetCaptchaForEmail(email, captcha string) error {
	key := GenerateCaptcha(email) //生成与邮箱绑定的验证码key
	expire := 2 * time.Minute
	return Rdb.Set(ctx, key, captcha, expire).Err()
	//向Redis写入key-value，并设置过期时间
	//Set返回*StatusCmd，通过Err()获取错误
}

// 用于校验用户输入的验证码是否正确
func CheckCaptchaForEmail(email, userInput string) (bool, error) {
	key := GenerateCaptcha(email)

	//从Redis中获取存储的验证码
	storedCaptcha, err := Rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			//表示key不存在（过期或从未设置），不算系统错误
			return false, nil
		}

		return false, err
	}

	//忽略大小写比较验证码
	if strings.EqualFold(storedCaptcha, userInput) {

		// 验证成功后删除 key，防止重复使用
		if err := Rdb.Del(ctx, key).Err(); err != nil {
			//即使删除失败，也不影响验证结果
		} else {
			//删除成功
		}
		return true, nil
	}

	//验证码不匹配
	return false, nil
}
