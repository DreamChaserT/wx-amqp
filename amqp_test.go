package wxamqp

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestWxAmqp(t *testing.T) {
	instance := NewWxAmqp("guest", "guest", "127.0.0.1", 5672)
	channel := instance.GetChannel()

	go func() {
		channel2 := instance.GetChannel()
		deliveryChan := channel2.AddConsumer("testQueue", "testConsumer", 1)

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

		time.Sleep(time.Second * 10)
		channel2.RemoveConsumer("testConsumer")
		channel2.Close()
	}()

	for i := 0; i < 20; i++ {
		err := channel.SendToQueue("testQueue", "hello test from client "+strconv.Itoa(i))
		fmt.Println(err)
		i++
		time.Sleep(time.Second)
	}
	instance.DisConnect()
	time.Sleep(time.Second * 10)
}
