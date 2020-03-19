package main

import (
	"fmt"
	// "strings"
	// snmp_serv "../internal/snmp_handler"
	// amqp "Snmp-Golang/internal/amqp_handler"

	amqp_serv "Snmp-Golang/internal/amqp_services"
	common_serv "Snmp-Golang/internal/common_services"
	model "Snmp-Golang/internal/models"
)


func main() {

	// snmpbulkget -v2c -Cn0 -Cr5 -c public 190.169.1.5 .1.3.6.1.4.1.119.2.3.69.501.7
	// model.LogSave("test1 test2 test3", "")

	commonParams := model.GetCommonInitParam()

	sendParams := model.SnmpSendParams{
		Ip:  "192.168.2.184",
		Oid: ".1.3.6.1",
		// Oid:".1.3.6.1.2.1.1.9.1",
		Community: "public",
		Port:      "161",
		DeviceId:  "TestDeviceId-777-999",
		SelCount:  0,
	}

	// queueAlias   = ""
	// amqpUrl      = "amqp://tryettt"
	// saveApiUrl   = model.SAVE_API_URL
	// amqpFuncType = model.AMQP_FUNC_TYPE


	//amqp.GetMessagesListStart(amqpUrl, queueName, saveApiUrl, amqpFuncType)

	// amqp.RecevieMessagesListFromQueue(amqpUrl, queueName, saveApiUrl)

	//queueParam, errOpen := amqp.RabbitQueueOpen(amqpUrl, queueAlias)
	//if errOpen != nil {
	//	amqp.FailOnError(errOpen, "RabbitQueueOpen - FATAL ERROR")
	//}
	//
	//newQueueName := queueParam.Name
	//msg, _ , err := queueParam.Channel.Get(newQueueName, true)
	//if err != nil {
	//	amqp.WarnOnError(err, "Not Get message in Queues - WARN ERROR")
	//}
	//
	//// _ , _ , e := queueParam.Channel.Get(newQueueName, true)
	//
	//item    := string(msg.Body)
	//message := strings.Split(item, " ")
	//fmt.Println(message, item)
	//
	//m := "3915 192.168.2.184 .1.3.6.1.2.1.1.9.1 161 SNMP 0"
	//errRetry := amqp.AmqpProducer(m)
	//fmt.Println("Producer error", errRetry)


	common_serv.GetQueueMessagesInizialize(commonParams, sendParams)

	// queueInitTest(amqpUrl, saveApiUrl, queueAlias, amqpFuncType)


}


func queueInitTest(amqpUrl, saveApiUrl, queueAlias, funcType string) {

	queueInit, errOpen := amqp_serv.RabbitQueueInit(amqpUrl, queueAlias)
	if errOpen != nil {
		model.FailOnError(errOpen, "RabbitQueueOpen - FATAL ERROR")
	}

	defer queueInit.Connect.Close()
	defer queueInit.Channel.Close()

	newQueueName := queueInit.Name
	msg, statOk , err := queueInit.Channel.Get(newQueueName, true)
	if err != nil {
		model.WarnOnError(err, "Not Get message in Queues - WARN ERROR")
	}

	if statOk {
		item := string(msg.Body)
		//	message := strings.Split(item, " ")
		// fmt.Println(queueInit)
		fmt.Println(item)
	} else {
		fmt.Println("Not message", statOk)
	}

}
