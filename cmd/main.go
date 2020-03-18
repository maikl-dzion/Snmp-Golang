package main

import (
	amqp "Snmp-Golang/internal/amqp_handler"
	model "Snmp-Golang/internal/models"
)

func main() {

	// snmpbulkget -v2c -Cn0 -Cr5 -c public 190.169.1.5 .1.3.6.1.4.1.119.2.3.69.501.7

	queueName  := model.QUEUE_NAME
	amqpUrl    := model.AMQP_API_URL
	saveApiUrl := model.SAVE_API_URL
	amqpFuncType := model.AMQP_FUNC_TYPE

	amqp.GetMessagesListStart(amqpUrl, queueName, saveApiUrl, amqpFuncType)
	// amqp.RecevieMessagesListFromQueue(amqpUrl, queueName, saveApiUrl)
}
