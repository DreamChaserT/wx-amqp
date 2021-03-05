package wxamqp

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func GetChannel() *AmqpChannel {
	instance := NewWxAmqp("guest", "guest", "127.0.0.1", 5672)
	return instance.GetChannel()
}

func TestGetChannel2(t *testing.T) {
	NewWxAmqpWithVhost("guest", "guest", "127.0.0.1", 5672, "tku")
	time.Sleep(time.Second * 10)
}

func TestAmqpChannel_DeclareQueue(t *testing.T) {
	channel := GetChannel()

	qName, err := channel.DeclareQueue("DeclareQueueNameTT", true, true)
	if nil != err {
		fmt.Println(err)
	} else {
		fmt.Println(qName)
	}

	qName, err = channel.DeclareQueue("DeclareQueueNameTF", true, false)
	if nil != err {
		fmt.Println(err)
	} else {
		fmt.Println(qName)
	}

	qName, err = channel.DeclareQueue("DeclareQueueNameFT", false, true)
	if nil != err {
		fmt.Println(err)
	} else {
		fmt.Println(qName)
	}

	qName, err = channel.DeclareQueue("DeclareQueueNameFF", false, false)
	if nil != err {
		fmt.Println(err)
	} else {
		fmt.Println(qName)
	}
}

func TestAmqpChannel_DeclareExchange(t *testing.T) {
	channel := GetChannel()

	err := channel.DeclareExchange("DeclareExchangeD", "direct", false, false)
	if nil != err {
		log.Println(err)
		t.FailNow()
	}

	err = channel.DeclareExchange("DeclareExchangeF", "fanout", false, false)
	if nil != err {
		log.Println(err)
		t.FailNow()
	}

	err = channel.DeclareExchange("DeclareExchangeT", "topic", false, false)
	if nil != err {
		log.Println(err)
		t.FailNow()
	}

	err = channel.DeclareExchange("DeclareExchangeH", "headers", false, false)
	if nil != err {
		log.Println(err)
		t.FailNow()
	}

	err = channel.DeclareExchange("DeclareExchangeTT", "direct", true, true)
	if nil != err {
		log.Println(err)
		t.FailNow()
	}

	err = channel.DeclareExchange("DeclareExchangeFT", "direct", false, true)
	if nil != err {
		log.Println(err)
		t.FailNow()
	}

	err = channel.DeclareExchange("DeclareExchangeTF", "direct", true, false)
	if nil != err {
		log.Println(err)
		t.FailNow()
	}
}

func TestBindAndUnBind(t *testing.T) {
	channel := GetChannel()
	channel.DeclareExchange("bindExchange", "direct", false, true)
	channel.DeclareQueue("bindQueue", false, true)
	err := channel.Bind("bindQueue", "*", "bindExchange")

	if nil != err {
		log.Println(err)
		t.FailNow()
	}

	err = channel.Unbind("bindQueue", "*", "bindExchange")
	if nil != err {
		log.Println(err)
		t.FailNow()
	}
}

func TestSendExchange(t *testing.T) {
	channel := GetChannel()
	channel.SendToQueue("PUSH_SERVICE_QUEUE", "hello")
	err := channel.SendToExchange("tb", "acc", "hello")
	log.Println(err)
	time.Sleep(time.Second * 10)
}
