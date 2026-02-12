package router

//用于在当前路由分组下注册一个POST类型的HTTP路由
//当客户端请求路径为/reconginze时，由image.RecognizeImage这个函数处理

import (
	"GopherAI/controller/image"

	"github.com/gin-gonic/gin"
)

func ImageRouter(r *gin.RouterGroup) {
	r.POST("/recognize", image.RecognizeImage)
}

//1.HTTP的常见方法：GET（获取数据），POST（提交数据）
//GET只能向服务器传递路径，没有请求体，而POST有（用于上传文件，图片等）
//2.该函数的作用是注册一个POST路由规则(根据路径和请求类型执行对应的方法)
//RouterGroup是路由分组对象
//3.最终执行示例
//router := gin.Default()
//imageGroup := router.Group("/image")
//ImageRouter(imageGroup)
