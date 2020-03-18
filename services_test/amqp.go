package main

import (
	// "fmt"
	// snmp_serv "../internal/snmp_handler"

	amqp "Snmp-Golang/internal/amqp_handler"
	model "Snmp-Golang/internal/models"
	"fmt"
	"strings"

	// "strings"
)


func main() {

	// snmpbulkget -v2c -Cn0 -Cr5 -c public 190.169.1.5 .1.3.6.1.4.1.119.2.3.69.501.7
	// model.LogSave("test1 test2 test3", "")

	queueName  := model.QUEUE_NAME
	amqpUrl    := model.AMQP_API_URL
	//saveApiUrl := model.SAVE_API_URL
	//amqpFuncType := model.AMQP_FUNC_TYPE

	//amqp.GetMessagesListStart(amqpUrl, queueName, saveApiUrl, amqpFuncType)

	// amqp.RecevieMessagesListFromQueue(amqpUrl, queueName, saveApiUrl)

	queueParam, errOpen := amqp.RabbitQueueOpen(amqpUrl, queueName)
	if errOpen != nil {
		amqp.FailOnError(errOpen, "RabbitQueueOpen - FATAL ERROR")
	}

	newQueueName := queueParam.Name
	msg, _ , err := queueParam.Channel.Get(newQueueName, true)
	if err != nil {
		amqp.WarnOnError(err, "Not Get message in Queues - WARN ERROR")
	}

	// _ , _ , e := queueParam.Channel.Get(newQueueName, true)

	item    := string(msg.Body)
	message := strings.Split(item, " ")
	fmt.Println(message, item)

	m := "3915 192.168.2.184 .1.3.6.1.2.1.1.9.1 161 SNMP 0"
	errRetry := amqp.AmqpProducer(m)
	fmt.Println("Producer error", errRetry)

}

