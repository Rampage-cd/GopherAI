package router

import (
	"GopherAI/controller/user"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.RouterGroup) {
	{
		//用户注册请求
		r.POST("/register", user.Register)
		//用户登录
		r.POST("/login", user.Login)
		//向指定邮箱发送验证码
		r.POST("/captcha", user.HandleCaptcha)
	}
}
