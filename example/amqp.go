package main

import (
	//"flag"
	//"log"

	"./consumer"
	"log"
)


const QUEUE_NAME   = "SNMP_QUEUE"
const SAVE_API_URL = "http://172.16.16.235:8080/data/save"
const AMQP_API_URL = "amqp://tester:12345@172.16.16.235:5672/"

func main() {

	queueParams := consumer.QueueInitParams{ false, false, false, false, QUEUE_NAME }


	channel , err := consumer.AmqpInitialize(AMQP_API_URL);
	// consumer.FailOnError(err, "RabbitMQ init error")
	// log.Println(channel, err)

	queue , err := consumer.QueueDeclareInit(channel, queueParams)
	consumer.FailOnError(err, "QueueDeclare init error")
	log.Println(queue, err)



	//ch, err := consumer.AmqpInitialize(AMQP_API_URL, QUEUE_NAME);
	//log.Println(ch, err)

	//var amqpUri string
	//
	//flag.StringVar(&amqpUri, "amqpUri", AMQP_API_URL, "AMQP connection uri")
	//flag.Parse()
	//
	//amqpConsumer := consumer.NewAmqpConsumer(amqpUri)
	//
	//messages     := amqpConsumer.ReceiveWithoutTimeout(
	//	                "events",
	//	                       []string{"efa1", "efa2"},
	//	                       QUEUE_NAME,
	//	                       consumer.QueueOptions{Durable: false, Delete: true, Exclusive: true})
	//
	//
	//for message := range messages {
	//	log.Println(message)
	//}

}