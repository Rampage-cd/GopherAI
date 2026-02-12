package router

import (
	"GopherAI/middleware/jwt"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default() //创建gin引擎（同时注册了Logger和Recovery两个中间件）
	enterRouter := r.Group("/api/v1")
	{
		//给用户请求注册路由
		RegisterUserRouter(enterRouter.Group("/user"))
	}
	//后续登录的接口需要jwt鉴权
	{
		AIGroup := enterRouter.Group("/AI")
		AIGroup.Use(jwt.Auth()) //绑定中间件，意味着这个Group下面的所有接口都必须经过JWT鉴权
		AIRouter(AIGroup)
	}

	{
		ImageGroup := enterRouter.Group("/image")
		ImageGroup.Use(jwt.Auth())
		ImageRouter(ImageGroup)
	}

	return r
}
