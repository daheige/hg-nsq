package main

import (
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
	//tcpNsqdAddrr = "127.0.0.1:4152"
	lookupdAddrr = "127.0.0.1:4161" //lookupd http地址
)

//声明一个结构体，实现HandleMessage接口方法（根据文档的要求）
type NsqHandler struct {
	//消息数
	msqCount int64
	//标识ID
	nsqHandlerID string
}

//实现 Handler接口上的HandleMessage方法
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
	conf := nsq.NewConfig()
	conf.ReadTimeout = 10 * time.Second
	conf.WriteTimeout = 10 * time.Second
	conf.HeartbeatInterval = 5 * time.Second //心跳检查

	//创造消费者，参数一时订阅的主题，参数二是使用的通道
	com, err := nsq.NewConsumer("test", "channel1", conf)
	if err != nil {
		fmt.Println(err)
	}

	//添加处理回调
	//com.AddHandler(&NsqHandler{nsqHandlerID: "One"}) //默认是单个goroutine处理消息

	//通过并发的方式消费
	com.AddConcurrentHandlers(&NsqHandler{nsqHandlerID: "One"}, 10)
	//连接对应的nsqd
	//err = com.ConnectToNSQD(tcpNsqdAddrr)

	//通过lookupd查询到nsqd节点后，连接到对应的nsqd
	err = com.ConnectToNSQLookupd(lookupdAddrr)
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

	//优雅的停止消费者
	com.Stop()

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
