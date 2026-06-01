package main

import (
	"GopherAI/common/aihelper"
	"GopherAI/common/mysql"
	"GopherAI/common/rabbitmq"
	"GopherAI/common/redis"
	"GopherAI/config"
	"GopherAI/dao/message"
	"GopherAI/router"
	"fmt"
	"log"
	"os"
	"os/signal"
	"net/http"
	"syscall"
	"time"
	"context"
)

func StartServer(addr string, port int) error {
	r := router.InitRouter()
	//服务器静态资源路径映射关系，这里目前不需要
	// r.Static(config.GetConfig().HttpFilePath, config.GetConfig().MusicFilePath)
	return r.Run(fmt.Sprintf("%s:%d", addr, port))
}

// 从数据库加载消息并初始化 AIHelperManager
func readDataFromDB() error {
	manager := aihelper.GetGlobalManager()
	// 从数据库读取所有消息
	msgs, err := message.GetAllMessages()
	if err != nil {
		return err
	}
	// 遍历数据库消息
	for i := range msgs {
		m := &msgs[i]
		//默认openai模型
		modelType := "1"
		config := make(map[string]interface{})

		// 创建对应的 AIHelper
		helper, err := manager.GetOrCreateAIHelper(m.UserName, m.SessionID, modelType, config)
		if err != nil {
			log.Printf("[readDataFromDB] failed to create helper for user=%s session=%s: %v", m.UserName, m.SessionID, err)
			continue
		}
		log.Println("readDataFromDB init:  ", helper.SessionID)
		// 添加消息到内存中(不开启存储功能)
		helper.AddMessage(m.Content, m.UserName, m.IsUser, false)
	}

	log.Println("AIHelperManager init success ")
	return nil
}

func main() {
	conf := config.GetConfig()
	host := conf.MainConfig.Host
	port := conf.MainConfig.Port
	//初始化mysql
	if err := mysql.InitMysql(); err != nil {
		log.Println("InitMysql error , " + err.Error())
		return
	}
	//初始化AIHelperManager
	readDataFromDB()

	//初始化redis
	if err := redis.Init(); err!=nil{
		log.Fatalf("redis init failed: %v",err)
	}
	log.Println("redis init success  ")
	if err := rabbitmq.InitRabbitMQ(); err!=nil{
		log.Fatalf("rabbitmq init failed: %v",err)
	}
	log.Println("rabbitmq init success  ")

	// err := StartServer(host, port) // 启动 HTTP 服务
	// if err != nil {
	// 	panic(err)
	// }
	//手动创建http.Server,方便后续进行优雅退出
	r := router.InitRouter()

	srv := &http.Server{
		Addr:		fmt.Sprintf("%s:%d",host,port),
		Handler:	r,
	}

	//开启独立协程启动Web服务
	go func(){
		log.Printf("GopherAI Server is running at http://%s:%d\n",host,port)
		if err := srv.ListenAndServe(); err!=nil && err != http.ErrServerClosed{
			log.Fatalf("HTTP Server listen error: %s\n",err)
		}
	}()

	//优雅退出机制
	quit := make(chan os.Signal,1)
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	
	//主协程阻塞，直到收到退出信号
	<-quit
	log.Printf("接收到停止信号,正在准备安全退出系统...")

	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()//这5秒钟用于处理还没有处理完成的旧请求

	if err := srv.Shutdown(ctx); err != nil{
		log.Fatalf("HTTP Server 强制关闭异常:",err)
	}

	//安全销毁消息队列
	rabbitmq.DestoryRabbitMQ()
	log.Println("RabbitMQ 连接已安全关闭")

	log.Println("GopherAI Server已优雅退出。")
}
