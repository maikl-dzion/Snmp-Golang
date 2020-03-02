package amqp_handler


import (
	model "../models"
	"fmt"
	mq "github.com/streadway/amqp"
	"log"
	"strings"
)

func MessagesRendering(messages <-chan mq.Delivery) {

	//port := "161"
	//oid  := "1.3.6.1.2.1.1.1.0";
	//ipAddress := "192.168.2.184"

	//sendParams := model.SnmpSendParams{
	//	"190.169.1.5",
	//	".1.3.6.1.4.1.119.2.3.69.501.7.1.1",
	//	"public",
	//	"161",
	//	0,
	//}

	for m := range messages {

		mess := string(m.Body)
		item := strings.Split(mess, " ")
		messageId := item[0];

		//ipAddress := item[1]
		//oid  := item[2]
		//port := "161"

		// snmp_handl.BulkRequestRun(sendParams)
		fmt.Println(messageId)

	}

}



func RecevieMessagesFromQueue() {

	connect, err := mq.Dial(model.AMQP_API_URL)
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer connect.Close()

	channel, err := connect.Channel()
	FailOnError(err, "Failed to open a channel")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		model.QUEUE_NAME,  // name
		false,      // durable
		false,    // delete when unused
		false,     // exclusive
		false,      // no-wait
		nil,          // arguments
	)

	FailOnError(err, "Failed to declare a queue")

	messages, err := channel.Consume(
		queue.Name,    // queue
		"",    // consumer
		true,   // auto-ack
		false, // exclusive
		false,  // no-local
		false,  // no-wait
		nil,      // args
	)

	FailOnError(err, "Failed to register a consumer")


	ch := 0

	for m := range messages {

		mess := string(m.Body)
		item := strings.Split(mess, " ")
		messageId := item[0];

		//ipAddress := item[1]
		//oid  := item[2]
		//port := "161"

		// snmp_handl.BulkRequestRun(sendParams)

		fmt.Println("I:", ch, "MessageId:", messageId)
		ch++

	}


	// fmt.Println(messages)

	//forever := make(chan bool)
	//
	//
	//for a := 0; a < 10; a++ {
	//	go MessagesRendering(messages)
	//}
	//
	//
	//log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	//
	//
	//<-forever

}



func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
