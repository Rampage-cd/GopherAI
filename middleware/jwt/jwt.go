package jwt

import (
	"GopherAI/common/code"
	"GopherAI/controller"
	"GopherAI/utils/myjwt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := new(controller.Response) //创建一个统一的响应结构体实例
		var token string
		authHeader := c.GetHeader("Authorization")
		//从HTTP Header中读取Authorization字段

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
			//去掉Bearer 前缀，提取真正的token字符串
		} else {
			token = c.Query("token") //兼容通过URL查询参数传递token的方式
		}

		if token == "" {
			c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidToken))
			//返回统一的无效的token错误码
			c.Abort() //终止当前请求的后续Handler执行
			return
		}

		log.Println("token is", token) //打印token到调试台（调试用）

		userName, ok := myjwt.ParseToken(token) //调用解析函数，返回从token中解析出的用户名
		if !ok {
			c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidToken))
			c.Abort()
			return
		}

		//如果解析成功：将userName存入Gin的上下文中
		//作用：后续Handler可以通过c.Get("userName")获取当前登录用户
		c.Set("userName", userName)

		//调用Next()，表示当前中间件逻辑结束，继续执行后续的Handler或中间件
		c.Next()

	}
}
