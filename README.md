# WX AMQP
## 增加特性
- 自动重连
- 简化操作
## 功能列表
- 获取包
```
go get -u github.com/DreamChaserT/wx-amqp
```
- 连接服务器
```
instance := wxamqp.NewWxAmqp("guest", "guest", "127.0.0.1", 5672)
instance := wxamqp.NewWxAmqpWithVhost("guest", "guest", "127.0.0.1", 5672,"vhost")
```
- 获取连接通道
```
channel := instance.GetChannel()
```
- 创建消息队列
```
qName, err := channel.DeclareQueue("队列名称", true, false)
```
- 创建交换机
```
err := channel.DeclareExchange("交换机名称", "direct", false, false)
```
- 绑定队列至交换机
```
err := channel.Bind("队列名称", "routingKey", "交换机名称")
```
- 解绑队列与交换机
```
err := channel.Unbind("队列名称", "routingKey", "交换机名称")
```
- 发送消息至指定队列
```
err := channel.SendToQueue("队列名称", "msg content")
```
- 发送消息至指定交换机
```
err := channel.SendToExchange("交换机名称","routingKey", "msg content")
```
- 监听队列,获取数据(手动ACK)
```
# consumerId 必须唯一
deliveryChan := channel.AddConsumer("队列名称", "consumerId", 1)
go func() {
    for {
        data, ok := <-deliveryChan
        if !ok {
            break
        }
        body := string(data.Body)
        fmt.Println(body)
        data.Ack(false)
    }
}()
```
- 取消队列监听
```
channel.RemoveConsumer("consumerId")
```
- 关闭连接通道
```
channel.Close()
```
- 断开服务器
```
instance.DisConnect()
```