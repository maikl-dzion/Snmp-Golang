package main

import (
	//"flag"
	//"log"

	"./consumer"
	_ "github.com/streadway/amqp"
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
	log.Println(queue)


	consumer.ReceiveQueueMessages(channel, queue.Name)

}