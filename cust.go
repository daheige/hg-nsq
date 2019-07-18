package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
)

var (
	//nsqd的地址，使用了tcp监听的端口
	// tcpNsqdAddrr = "127.0.0.1:4150" //消费第一个insert
	tcpNsqdAddrr = "127.0.0.1:4152"
)

//声明一个结构体，实现HandleMessage接口方法（根据文档的要求）
type NsqHandler struct {
	//消息数
	msqCount int64
	//标识ID
	nsqHandlerID string
}

//实现HandleMessage方法
//message是接收到的消息
func (s *NsqHandler) HandleMessage(message *nsq.Message) error {
	//没收到一条消息+1
	s.msqCount++
	//打印输出信息和ID
	fmt.Println(s.msqCount, s.nsqHandlerID)
	//打印消息的一些基本信息
	fmt.Printf("msg.Timestamp=%v, msg.nsqaddress=%s,msg.body=%s \n", time.Unix(0, message.Timestamp).Format("2006-01-02 03:04:05"), message.NSQDAddress, string(message.Body))
	return nil
}

func main() {
	//初始化配置
	config := nsq.NewConfig()
	//创造消费者，参数一时订阅的主题，参数二是使用的通道
	com, err := nsq.NewConsumer("Insert", "channel1", config)
	if err != nil {
		fmt.Println(err)
	}
	//添加处理回调
	com.AddHandler(&NsqHandler{nsqHandlerID: "One"})
	//连接对应的nsqd
	err = com.ConnectToNSQD(tcpNsqdAddrr)
	if err != nil {
		fmt.Println(err)
	}

	//平滑重启
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	//window signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	sig := <-ch

	log.Println("exit signal: ", sig.String())
	// Create a deadline to wait for.
	_, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	log.Println("shutting down")
}

/**
sg.Timestamp=2019-07-18 10:37:21, msg.nsqaddress=127.0.0.1:4152,msg.body=hello:94
96 One
msg.Timestamp=2019-07-18 10:37:21, msg.nsqaddress=127.0.0.1:4152,msg.body=hello:95
97 One
msg.Timestamp=2019-07-18 10:37:21, msg.nsqaddress=127.0.0.1:4152,msg.body=hello:96
98 One
msg.Timestamp=2019-07-18 10:37:21, msg.nsqaddress=127.0.0.1:4152,msg.body=hello:97
99 One
msg.Timestamp=2019-07-18 10:37:21, msg.nsqaddress=127.0.0.1:4152,msg.body=hello:98
100 One
msg.Timestamp=2019-07-18 10:37:21, msg.nsqaddress=127.0.0.1:4152,msg.body=hello:99
*/
