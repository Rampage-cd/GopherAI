package router

import (
	"GopherAI/controller/session"

	"github.com/gin-gonic/gin"
)

func AIRouter(r *gin.RouterGroup) {
	//聊天相关接口
	{
		//获取用户所有会话
		r.GET("/chat/sessions", session.GetUserSessionsByUserName)
		//创建新会话并发送消息
		r.POST("/chat/send-new-session", session.CreateSessionAndSendMessage)
		//在已创建的会话上发送消息
		r.POST("/chat/send", session.ChatSend)
		//获取一次会话中的所有信息
		r.POST("/chat/history", session.ChatHistory)
		//创建流式输出会话
		r.POST("/chat/send-stream-new-session", session.CreateStreamSessionAndSendMessage)
		r.POST("/chat/send-stream", session.ChatStreamSend)
	}
}
