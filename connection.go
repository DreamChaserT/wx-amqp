package wxamqp

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

// AmqpConnection connection
type AmqpConnection struct {
	// amqp 连接信息
	username string
	password string
	host     string
	port     int
	vhost    string

	// 真实的连接
	c *amqp.Connection
	// 自动连接
	autoConnect bool
}

// 执行连接
func (a *AmqpConnection) connect() error {
	connection, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d%s", a.username, a.password, a.host, a.port, a.vhost))
	if nil != err {
		return err
	}
	a.c = connection
	// 启动监控
	a.monitorConnect()
	return nil
}

// 监控连接并重连
func (a *AmqpConnection) monitorConnect() {
	go func(instance *AmqpConnection) {
		for {
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
					break
				}
			} else {
				// 连接已停止,终止监控任务
				break
			}
		}
	}(a)
}

// Disconnect 断开连接
func (a *AmqpConnection) Disconnect() {
	a.autoConnect = false
	if nil != a.c {
		a.c.Close()
	}
}

// NewAmqpConnection new
func NewAmqpConnection(username string, password string, host string, port int, vhost string) *AmqpConnection {
	instance := &AmqpConnection{
		username:    username,
		password:    password,
		host:        host,
		port:        port,
		vhost:       vhost,
		autoConnect: true,
	}
	instance.monitorConnect()
	return instance
}
