package main

import (

	rabbitmq "../internal/amqp_handler"
	model "../internal/models"
)


func main() {

	queueName := model.QUEUE_NAME
	amqpUrl   := model.AMQP_API_URL

	rabbitmq.RecevieMessagesListFromQueue(amqpUrl, queueName)

}
