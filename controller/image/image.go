package image

import (
	"GopherAI/common/code"
	"GopherAI/controller"
	"GopherAI/service/image"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	//gin作为Web框架，作用如下
	//解析HTTP请求，路由匹配，参数获取，返回JSON响应
)

type (
	RecognizeImageResponse struct {
		ClassName           string `json:"class_name,omitempty"` // AI回答
		controller.Response        //匿名嵌套，可以直接使用Response结构体中的字段和方法
	}
)

func RecognizeImage(c *gin.Context) { //*gin.Context是Gin的Handler函数
	//每个HTTP请求，Gin都会创建一个*gin.Context
	res := new(RecognizeImageResponse)
	file, err := c.FormFile("image")
	//该函数解析了multipart/form-data，并从请求中找到名为"image"的文件，返回*multipart.FileHeader
	if err != nil {
		log.Println("FormFile fail ", err)
		c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidParams))
		return //设置HTTP状态码，自动把后面的data序列化为JSON
	}

	className, err := image.RecognizeImage(file)
	if err != nil {
		log.Println("RecognizeImage fail ", err)
		c.JSON(http.StatusOK, res.CodeOf(code.CodeServerBusy))
		return
	}

	res.Success()
	res.ClassName = className
	c.JSON(http.StatusOK, res)
}
