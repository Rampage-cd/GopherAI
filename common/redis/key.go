package redis

import (
	"GopherAI/config"
	"fmt"
)

func GenerateCaptcha(email string) string {
	return fmt.Sprintf(config.DefaultRedisKeyConfig.CaptchaPrefix, email)
	//返回的结果是"captcha:%s",%s中的内容为email
}
