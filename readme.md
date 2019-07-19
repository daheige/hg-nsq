# nsq消息队列
    nsq学习

# 下载安装
    cd /usr/local
    sudo wget https://s3.amazonaws.com/bitly-downloads/nsq/nsq-1.1.0.linux-amd64.go1.10.3.tar.gz
    sudo tar zxvf nsq-1.1.0.linux-amd64.go1.10.3.tar.gz
    sudo tar mv nsq-1.1.0.linux-amd64.go1.10.3 nsq
    sudo chown -R $USER nsq

    vim ~/.bashrc
        export NSQ_HOME=/usr/local/nsq
        export PATH=$NSQ_HOME/bin:$PATH
    source ~/.bashrc
# 运行nsq
    1、启动nsqlookupd
        $ nsqlookupd
        [nsqlookupd] 2019/07/18 22:07:53.028370 INFO: nsqlookupd v1.1.0 (built w/go1.10.3)
        [nsqlookupd] 2019/07/18 22:07:53.028664 INFO: HTTP: listening on [::]:4161
        [nsqlookupd] 2019/07/18 22:07:53.028676 INFO: TCP: listening on [::]:4160

    2、启动两个nsqd

        mkdir -p /usr/local/nsq/data && nsqd --lookupd-tcp-address=127.0.0.1:4160 -tcp-address=127.0.0.1:4152 -http-address=127.0.0.1:4153 -broadcast-address=127.0.0.1 --data-path=/usr/local/nsq/data

        mkdir -p /usr/local/nsq/data2 && nsqd --lookupd-tcp-address=127.0.0.1:4160 -tcp-address=127.0.0.1:4150 -http-address=127.0.0.1:4151 -broadcast-address=127.0.0.1 --data-path=/usr/local/nsq/data2

    3、启动nqsadmin，连接到nsqlookupd
        nsqadmin --lookupd-http-address=127.0.0.1:4161
        访问http://localhost:4171 查看nsqdadmin

# 关于nsqdadmin
    我们先来说明一下这个后台里面的一些内容，因为我们的NSQ所使用的是经典的pub/sub模式（发布/订阅，典型的生产者/消费者模式），我们可以先发布一个主题到NSQ，然后所有订阅的服务器就会异步的从这里读取主题的内容：

    Topic(左上角)：发布的主题名字

    NSQd Host：Nsq主机服务地址

    Channel：消息通道

    NSQd Host：Nsq主机服务地址

    Depth：消息积压量

    In-flight：已经投递但是还未消费掉的消息

    Deferred：没有消费掉的延时消息

    Messages：服务器启动之后，总共接收到的消息量

    Connections：通道里面客户端的订阅数

    TimeOut：超时时间内没有被响应的消息数

    Memory + Disk：储存在内存和硬盘中总共的消息数

# 如何在代码中使用nsq
    如何在代码中发布主题内容，然后通过订阅某主题去异步读取消息

    使用官方提供的下载地址：

    go get github.com/nsqio/go-nsq

    先创建一个主题，并且发布100条消息
    代码参考：production.go
    发送消息
    go run production.go
    访问http://localhost:4171 就可以看到刚才发送后的消息topic
    可以看到Nsqd接收到了100条信息，100条信息都储存在内存中，没有被消化

# 创建nsq消费者
    现在没有任何服务订阅了我们的主题，所以主题的消息都没有被消化，那我们创建一个消费者去订阅我们的主题：
    参考cust.go

# 开始消费
    go run cust.go
    msg.Timestamp=2019-07-18 11:01:47, msg.nsqaddress=127.0.0.1:4152,msg.body=hello:97 
    99 One
    msg.Timestamp=2019-07-18 11:01:47, msg.nsqaddress=127.0.0.1:4152,msg.body=hello:98 
    100 One
    msg.Timestamp=2019-07-18 11:01:47, msg.nsqaddress=127.0.0.1:4152,msg.body=hello:99
    再次访问http://localhost:4171  之前挤压的100条信息，都被我们的订阅者消化掉了，也就是读取了

    所以我们的订阅者（可以有多个）如果提前订阅主题的话，只要对应的主题有发布新内容，就可以马上异步读取

    查看消费次数： http://localhost:4171/counter
    查看nsqd 4153 http://localhost:4171/nodes/127.0.0.1:4153

    可以启动多个消费者，这样的话，处理速度就会加快

# 官网文档
    https://nsq.io/overview/quick_start.html

# 优雅地退出发送和消费者模式
    参考production.go和cust.go 接收信号量处理方式
    
# 参考文档
    https://segmentfault.com/a/1190000009194607
    
    https://blog.csdn.net/sd653159/article/details/83624661

    https://github.com/nsqio/nsq

    https://github.com/nsqio/go-nsq
