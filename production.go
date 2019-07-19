package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

	//平滑重启
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	//window signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	//测试发送100w消息
	nums := 100 * 10000
	for i := 0; i < nums; i++ {
		//监听是否有退出信号
		select {
		case sig := <-ch: //接收到停止信号，就优雅的退出发送
			signal.Stop(ch) //停止接收信号

			log.Println("exit signal: ", sig.String())
			//优雅的停止发送
			//Stop initiates a graceful stop of the Producer (permanent)
			tPro.Stop()
			goto exit //如果是退出函数可以写return
		default:
			log.Println("msg sending...")
		}

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

exit:
	log.Println("production will exit")

}
