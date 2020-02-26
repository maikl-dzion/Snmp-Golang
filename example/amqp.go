package main

import (
	"flag"
	"log"

	"./consumer"
)


const QUEUE_NAME   = "SNMP_QUEUE"
const SAVE_API_URL = "http://172.16.16.235:8080/data/save"
const AMQP_API_URL = "amqp://tester:12345@172.16.16.235:5672/"

func main() {

	var amqpUri string

	flag.StringVar(&amqpUri, "amqpUri", AMQP_API_URL, "AMQP connection uri")
	flag.Parse()

	amqpConsumer := consumer.NewAmqpConsumer(amqpUri)

	messages     := amqpConsumer.ReceiveWithoutTimeout(
		                "events",
		                       []string{"efa1", "efa2"},
		                       QUEUE_NAME,
		                       consumer.QueueOptions{Durable: false, Delete: true, Exclusive: true})


	for message := range messages {
		log.Println(message)
	}

}