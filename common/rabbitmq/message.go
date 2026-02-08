package rabbitmq

import (
	"GopherAI/dao/message"
	"GopherAI/model"
	"encoding/json"

	"github.com/streadway/amqp"
)

type MessageMQParam struct {
	SessionID string `json:"session_id"` //会话ID
	Content   string `json:"content"`    //消息内容
	UserName  string `json:"user_name"`  //用户名
	IsUser    bool   `json:"is_user"`    //是否为用户消息
}

// 将消息数据序列化为JSON
// 用于投递到RabbitMQ
func GenerateMessageMQParam(sessionID string, content string, userName string, IsUser bool) []byte {
	param := MessageMQParam{
		SessionID: sessionID,
		Content:   content,
		UserName:  userName,
		IsUser:    IsUser,
	}
	data, _ := json.Marshal(param)
	return data
}

// RabbitMQ消费端的业务处理函数
func MQMessage(msg *amqp.Delivery) error {
	var param MessageMQParam
	err := json.Unmarshal(msg.Body, &param) //反序列化消息体
	if err != nil {
		return err
	}

	//转化为数据库模型
	newMsg := &model.Message{
		SessionID: param.SessionID,
		Content:   param.Content,
		UserName:  param.UserName,
		IsUser:    param.IsUser,
	}

	//消费者异步插入到数据库中
	message.CreateMessage(newMsg)
	return nil
}
