package public

// WxAmqp 主要功能类,用于操作消息队列
type WxAmqp struct {
	// amqp 连接信息
	username string
	password string
	host     string
	port     int
	c        *AmqpConnection
}

// NewWxAmqp 创建一个新的连接实例
func NewWxAmqp(username string, password string, host string, port int) *WxAmqp {
	return &WxAmqp{
		username: username,
		password: password,
		host:     host,
		port:     port,
		c:        NewAmqpConnection(username, password, host, port),
	}
}

// DisConnect 断开消息队列的连接
func (a *WxAmqp) DisConnect() {
	if nil != a.c {
		a.c.Disconnect()
	}
}

// GetChannel 获取连接通道
func (a *WxAmqp) GetChannel() *AmqpChannel {
	return NewAmqpChannel(a.c)
}
