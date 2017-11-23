package main

import (
	"github.com/nsqio/go-nsq"
	"log"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer("write_test", "ch", config)
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Printf("Received a message: %v", message)
		wg.Done()
		return nil
	}))

	if err := q.ConnectToNSQD("127.0.0.1:4150"); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
