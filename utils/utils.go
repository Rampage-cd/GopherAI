package utils

import (
	"GopherAI/model"
	"crypto/md5"   //提供MD5哈希算法的实现
	"encoding/hex" //用于将二进制数据编码成十六进制字符串
	"math/rand"    //提供伪随机数生成器
	"strconv"      //用于字符串与基础类型之间的转换
	"time"

	"github.com/cloudwego/eino/schema" //AI对话中的“消息抽象结构”
	"github.com/google/uuid"           //生成RFC 4122规范的UUID
)

// 生成一个指定（num）长度的纯数字随机字符串
func GetRandomNumbers(num int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//使用当前时间的纳秒作为种子

	code := ""
	//用于拼接最终生成的随机数字字符串
	for i := 0; i < num; i++ {
		// 0~9随机数
		digit := r.Intn(10)
		code += strconv.Itoa(digit)
	}
	return code
}

// MD5 MD5加密
func MD5(str string) string {
	m := md5.New() //返回hash.Hash

	//向hash中写入数据
	//Write可以多次调用
	m.Write([]byte(str))

	//Sum返回最终的二进制哈希结果（[]byte）
	//hex.EncodeToString将其编码成十六进制字符串
	return hex.EncodeToString(m.Sum(nil))
}

// 生成一个UUID v4字符串：基于随机数，冲突概率极低，不依赖中心化ID生成器
func GenerateUUID() string {
	return uuid.New().String() //String表示转化为标准字符串表示
}

// 将 schema 消息转换为数据库可存储的格式
func ConvertToModelMessage(sessionID string, userName string, msg *schema.Message) *model.Message {
	return &model.Message{
		SessionID: sessionID,   //会话ID，用于关联一轮对话
		UserName:  userName,    //用户名，用于区分不同用户
		Content:   msg.Content, //消息内容
	}
}

// 将数据库消息转换为 schema 消息（供 AI 使用）
func ConvertToSchemaMessages(msgs []*model.Message) []*schema.Message {
	schemaMsgs := make([]*schema.Message, 0, len(msgs))
	for _, m := range msgs {
		role := schema.Assistant //默认角色
		if m.IsUser {
			role = schema.User //用户角色
		}
		schemaMsgs = append(schemaMsgs, &schema.Message{
			Role:    role, //用于判断内容为用户输入，还是AI输出
			Content: m.Content,
		})
	}
	return schemaMsgs
}
