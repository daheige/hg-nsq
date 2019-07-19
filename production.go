package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/nsqio/go-nsq"
)

var (
	//nsqd的地址，使用了tcp监听的端口
	//这里是直接连接到nsqd,两个nsqd
	// tcpNsqdAddrr = "127.0.0.1:4150" //nsqd连接的tcp地址
	tcpNsqdAddrr = "127.0.0.1:4152" //nsqd连接的tcp地址
)

//先创建一个主题，并且发布100条消息
func main() {
	//初始化配置
	conf := nsq.NewConfig()
	conf.ReadTimeout = 10 * time.Second
	conf.WriteTimeout = 10 * time.Second
	conf.HeartbeatInterval = 5 * time.Second //心跳检查

	//创建生产者
	tPro, err := nsq.NewProducer(tcpNsqdAddrr, conf)
	if err != nil {
		fmt.Println("new producer err:", err)
	}

	//测试发送100w消息
	nums := 100 * 10000
	for i := 0; i < nums; i++ {
		//主题
		topic := "test"
		//主题内容
		tCommand := "hello:" + strconv.Itoa(i)
		//发布消息
		err = tPro.Publish(topic, []byte(tCommand))
		if err != nil {
			fmt.Println("publis msg error: ", err)
			continue
		}

		fmt.Println("current index: ", i)
		fmt.Println("send msg success")
	}
}
