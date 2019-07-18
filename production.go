package main

import (
	"fmt"
	"strconv"

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
	config := nsq.NewConfig()
	for i := 0; i < 100; i++ {
		//创建100个生产者
		tPro, err := nsq.NewProducer(tcpNsqdAddrr, config)
		if err != nil {
			fmt.Println("new producer err:", err)
			continue
		}

		//主题
		topic := "Insert"
		//主题内容
		tCommand := "hello:" + strconv.Itoa(i)
		//发布消息
		err = tPro.Publish(topic, []byte(tCommand))
		if err != nil {
			fmt.Println("publis msg error: ", err)
			continue
		}

		fmt.Println("send msg success")
	}
}
