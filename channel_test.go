package wxamqp

import (
	"fmt"
	"log"
	"testing"
)

func GetChannle() *AmqpChannel {
	instance := NewWxAmqp("guest", "guest", "127.0.0.1", 5672)
	return instance.GetChannel()
}

func TestAmqpChannel_DeclareQueue(t *testing.T) {
	channel := GetChannle()

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
	channel := GetChannle()

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
	channel := GetChannle()
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
