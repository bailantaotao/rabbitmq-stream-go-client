package main

import (
	"bufio"
	"fmt"
	"github.com/Azure/go-amqp"
	"github.com/google/uuid"
	"github.com/gsantomaggio/go-stream-client/pkg/stream"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	fmt.Println("RabbitMQ golang streaming client")
	fmt.Println("Connecting ...")
	var client = stream.NewStreamingClient()                                  // create Client Struct
	err := client.Connect("rabbitmq-stream://guest:guest@localhost:5551/%2f") // Connect
	if err != nil {
		fmt.Printf("error: %s", err)
		return
	}
	fmt.Println("Connected!")
	streamName := uuid.New().String()
	_, err = client.CreateStream(streamName) // Create the streaming queue
	if err != nil {
		fmt.Printf("error: %s", err)
		return
	}
	producers := 10
	iterations := 2000
	messages := 100
	count := 0
	consumer, err := client.NewConsumer(streamName, func(subscriberId byte, message *amqp.Message) {
		count++
		if count%(messages*iterations*producers) == 0 {
			fmt.Printf("Consumer id: %d - body: %s - Total messages consumed:%d \n", subscriberId, message.Data, count)
		}
	})
	if err != nil {
		fmt.Printf("error: %s", err)
		return
	}

	fmt.Printf("starting with %d producers, %d messages and  %d iterations \n", producers, iterations, messages)
	var arr []*amqp.Message // amqp 1.0 message from https://github.com/Azure/go-amqp
	for z := 0; z < messages; z++ {
		arr = append(arr, amqp.NewMessage([]byte("hello world_"+strconv.Itoa(z))))
	}
	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < producers; i++ {
		wg.Add(1)
		producer, err := client.NewProducer(streamName) // Get a new subscribe to publish the messages
		if err != nil {
			fmt.Printf("error: %s", err)
			return
		}
		go func(id int, producer *stream.Producer, wg *sync.WaitGroup) {
			defer wg.Done()
			for z := 0; z < iterations; z++ {
				_, err = producer.BatchPublish(nil, arr) // batch send
				if err != nil {
					fmt.Printf("error: %s", err)
					return
				}
			}
		}(i, producer, &wg)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("%s  to publish %d messages\n", elapsed, producers*iterations*messages)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press any key to finish ")
	_, _ = reader.ReadString('\n')
	_ = consumer.UnSubscribe()
	fmt.Println("Closing all producers ")
	err = client.CloseAllProducers()
	if err != nil {
		fmt.Printf("error: %s", err)
		return
	}
	_, err = client.DeleteStream(streamName) // Remove the streaming queue and the data
	if err != nil {
		fmt.Printf("error: %s", err)
		return
	}
	fmt.Println("Bye bye")
}
