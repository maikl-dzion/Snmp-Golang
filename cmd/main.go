package main

import (
	amqp "github.com/maikl-dzion/Snmp-Golang/internal/amqp_handler"
	model "github.com/maikl-dzion/Snmp-Golang/internal/models"
)


func main() {

	queueName := model.QUEUE_NAME
	amqpUrl   := model.AMQP_API_URL

	amqp.RecevieMessagesListFromQueue(amqpUrl, queueName)

	// fmt.Println(queueName, amqpUrl)

}
