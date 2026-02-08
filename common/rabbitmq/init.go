package rabbitmq

var (
	RMQMessage *RabbitMQ
)

func InitRabbitMQ() {
	//创建MQ并启动消费者
	//无论调用多少次NewWorkRabbitMQ，只会创建一次连接
	//不同队列公用一个连接，可以保持不同队列消费消息的顺序

	RMQMessage = NewWorkRabbitMQ("Message")
	go RMQMessage.Consume(MQMessage)

}

// 销毁消息队列
func DestoryRabbitMQ() {
	RMQMessage.Destroy()
}
