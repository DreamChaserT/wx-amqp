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
```
- 获取连接通道
```
channel := instance.GetChannel()
```
- 发送消息至指定队列
```
err := channel.SendToQueue("queue_name", "msg content")
```
- 监听队列,获取数据(手动ACK)
```
# consumerId 必须唯一
deliveryChan := channel.AddConsumer("queue_name", "consumerId", 1)
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