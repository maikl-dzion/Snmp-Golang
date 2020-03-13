package main

import (
	amqp "../internal/amqp_handler"
	model "../internal/models"
)

// snmpbulkget -v2c -Cn0 -Cr5 -c public 190.169.1.5 .1.3.6.1.4.1.119.2.3.69.501.7

func main() {

	queueName  := model.QUEUE_NAME
	amqpUrl    := model.AMQP_API_URL
	saveApiUrl := model.SAVE_API_URL

	amqp.RecevieMessagesListFromQueue(amqpUrl, queueName, saveApiUrl)

}
