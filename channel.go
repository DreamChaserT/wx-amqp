package wxamqp

import (
	"errors"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

// ErrDisconnect channel 断线异常
var ErrDisconnect = errors.New("channel 断线")

// AmqpChannel channel
type AmqpChannel struct {
	connection **amqp.Connection
	// 真实的连接
	c *amqp.Channel
	// 自动连接
	autoConnect bool
	// 消费者
	consumers sync.Map
}

// 执行连接
func (a *AmqpChannel) connect() error {
	if nil == *(a.connection) {
		return errors.New("connection为空")
	}
	channel, err := (*a.connection).Channel()
	if nil != err {
		return err
	}
	a.c = channel
	// 启动监控
	a.monitorConnect()
	return nil
}

// 监控并重连
func (a *AmqpChannel) monitorConnect() {
	go func(instance *AmqpChannel) {
		for {

			if nil == (*a.connection) {
				continue
			}

			if (*a.connection).IsClosed() {
				// 外层连接已经关闭,停止监控
				break
			}

			// 连接已停止,终止监控任务
			if !a.autoConnect {
				break
			}

			if nil == instance.c {
				// 初次连接失败,执行连接
				err := instance.connect()
				if nil != err {
					// 首次连接失败,等待20ms
					time.Sleep(time.Duration(20) * time.Millisecond)
					continue
				} else {
					// 新监控已启动,停止
					break
				}
			}
			notifyChan := instance.c.NotifyClose(make(chan *amqp.Error))
			<-notifyChan
			// 连接断开,重置connection为nil
			instance.c = nil
			if instance.autoConnect {
				// 执行重连
				err := instance.connect()
				if nil != err {
					// 连接失败,等待20ms
					time.Sleep(time.Duration(20) * time.Millisecond)
					continue
				} else {
					// 新监控已启动,停止
					break
				}
			} else {
				// 连接已停止,终止监控任务
				break
			}
		}
	}(a)
}

// NewAmqpChannel new
func NewAmqpChannel(amqpConnection *AmqpConnection) *AmqpChannel {
	instance := &AmqpChannel{
		connection:  &amqpConnection.c,
		autoConnect: true,
	}
	instance.monitorConnect()
	for {
		if nil != instance.c {
			break
		}
	}
	return instance
}

// Close 关闭连接
func (a *AmqpChannel) Close() {
	if nil != a.c {
		a.autoConnect = false
		a.c.Close()
	}
}

// SendToExchange 发送数据至指定交换机
func (a *AmqpChannel) SendToExchange(exchangeName string, routingKey string, body string) error {
	if nil != a.c {
		return a.c.Publish(exchangeName, routingKey, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	}
	return ErrDisconnect
}

// SendToQueue 发送数据至指定队列
func (a *AmqpChannel) SendToQueue(queueName string, body string) error {
	return a.SendToExchange("", queueName, body)
}

// DeclareQueue 创建队列
func (a *AmqpChannel) DeclareQueue(queueName string, durable bool, autoDelete bool) (string, error) {
	if nil != a.c {
		q, err := a.c.QueueDeclare(queueName, durable, autoDelete, false, false, nil)
		if nil == err {
			return q.Name, nil
		}
		return "", err
	}
	return "", ErrDisconnect
}

// DeclareExchange 创建Exchange
// kind direct fanout topic headers
func (a *AmqpChannel) DeclareExchange(exchangeName, kind string, durable, autoDelete bool) error {
	if nil != a.c {
		return a.c.ExchangeDeclare(exchangeName, kind, durable, autoDelete, false, false, nil)
	}
	return ErrDisconnect
}

// Bind 绑定queue至exchange
func (a *AmqpChannel) Bind(queueName string, routingKey string, exchangeName string) error {
	if nil != a.c {
		return a.c.QueueBind(queueName, routingKey, exchangeName, false, nil)
	}
	return ErrDisconnect
}

// Unbind 解绑queue与exchange
func (a *AmqpChannel) Unbind(queueName string, routingKey string, exchangeName string) error {
	if nil != a.c {
		return a.c.QueueUnbind(queueName, routingKey, exchangeName, nil)
	}
	return ErrDisconnect
}

// AddConsumer 消费队列数据 手动ack
func (a *AmqpChannel) AddConsumer(queueName string, consumerName string, prefetch int) chan amqp.Delivery {
	deliveryChan := make(chan amqp.Delivery)
	a.consumers.Store(consumerName, "")
	go func(a *AmqpChannel, queueName string, deliveryChan chan amqp.Delivery) {
		for {
			_, ok := a.consumers.Load(consumerName)
			if !ok {
				// 取消监听
				close(deliveryChan)
				break
			} else {
				if nil != a.c {
					err := a.c.Qos(prefetch, 0, false)
					if nil != err {
						time.Sleep(time.Millisecond * 20)
						continue
					}
					deliveryChanInner, err := a.c.Consume(queueName, consumerName, false, false, false, false, nil)
					if nil != err {
						time.Sleep(time.Millisecond * 20)
						continue
					} else {
						for {
							msg, ok := <-deliveryChanInner
							if !ok {
								break
							} else {
								deliveryChan <- msg
							}
						}
					}
				}
			}
		}
	}(a, queueName, deliveryChan)
	return deliveryChan
}

// RemoveConsumer 移除消费者
func (a *AmqpChannel) RemoveConsumer(consumerName string) {
	a.consumers.Delete(consumerName)
	if nil != a.c {
		a.c.Cancel(consumerName, false)
	}
}
