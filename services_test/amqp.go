package main

import (
	// "fmt"
	// snmp_serv "../internal/snmp_handler"

	amqp "Snmp-Golang/internal/amqp_handler"
	model "Snmp-Golang/internal/models"
)


func main() {

	// model.LogSave("test1 test2 test3", "")

	queueName  := model.QUEUE_NAME
	amqpUrl    := model.AMQP_API_URL
	saveApiUrl := model.SAVE_API_URL
	amqpFuncType := model.AMQP_FUNC_TYPE

	amqp.GetMessagesListStart(amqpUrl, queueName, saveApiUrl, amqpFuncType)

	// amqp.RecevieMessagesListFromQueue(amqpUrl, queueName, saveApiUrl)

}
