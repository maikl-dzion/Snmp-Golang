package main

import (
	"../internal/handlers"
	"../internal/models"
	// "fmt"
	"github.com/streadway/amqp"
	"log"
	"strings"
	// "reflect"
)

func MessagesRender(messages <-chan amqp.Delivery) {

	port := "161"
	oid  := "1.3.6.1.2.1.1.1.0";
	ipAddress := "192.168.2.184"

    _count := 0

	for d := range messages {

		mess := string(d.Body)
		item := strings.Split(mess, " ")

		//port := "161"
		//oid  := item[2]
		//ipAddress := item[1]

		//if oid != "signal_in" && oid != "signal_out"

		messageId := item[0];
		handler.SnmpRequestInit(ipAddress, oid, port, messageId)

		log.Printf("Count: %d, MessageId : %s", _count, messageId)
		_count++

	}

}

// fmt.Println(reflect.TypeOf(mess))
// cd /home/dev/web/loader/log
// wc db.log
// sudo rm db.log


func main() {

	connect, err := amqp.Dial(model.AMQP_API_URL)
	handler.FailOnError(err, "Failed to connect to RabbitMQ")
	defer connect.Close()

	channel, err := connect.Channel()
	handler.FailOnError(err, "Failed to open a channel")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		model.QUEUE_NAME,  // name
		false,      // durable
		false,    // delete when unused
		false,     // exclusive
		false,      // no-wait
		nil,          // arguments
	)

	handler.FailOnError(err, "Failed to declare a queue")

	messages, err := channel.Consume(
		queue.Name,    // queue
		"",    // consumer
		true,   // auto-ack
		false, // exclusive
		false,  // no-local
		false,  // no-wait
		nil,      // args
	)

	handler.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	// fmt.Println(reflect.TypeOf(ch))

	//go func() {
		// go MessagesRender(messages)
	//}()


	//for a := 0; a < 150; a++ {
	//
	//	go MessagesRender(messages)
	//
	//}


	go MessagesRender(messages)

	go MessagesRender(messages)

	go MessagesRender(messages)

	go MessagesRender(messages)

	go MessagesRender(messages)

	go MessagesRender(messages)

	go MessagesRender(messages)


	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")


	<-forever
}