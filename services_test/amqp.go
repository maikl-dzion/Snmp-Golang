package main

import (
	// "fmt"
	// snmp_serv "../internal/snmp_handler"

	amqp "../internal/amqp_handler"
	model "../internal/models"
)


func main() {

	queueName  := model.QUEUE_NAME
	amqpUrl    := model.AMQP_API_URL
	saveApiUrl := model.SAVE_API_URL
	selectType := "get"

	// amqp.RecevieMessagesListFromQueue(amqpUrl, queueName, saveApiUrl)

	amqp.GetMessagesListStart(amqpUrl, queueName, saveApiUrl, selectType)


}



//func (r ResponseJsonItems) FomattedJson(d model.ResponseMessage) {
//	r = append(r, []byte(d))
//}