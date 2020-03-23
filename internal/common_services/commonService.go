package common_services

import (
	amqp_serv "Snmp-Golang/internal/amqp_services"
	model "Snmp-Golang/internal/models"
	snmp_serv "Snmp-Golang/internal/snmp_services"
	"fmt"
	"sync"
)


/////////////////////////////////
/////////////////////////////////
//    Common Amqp Services

func GetQueueMessagesInizialize(param model.CommonInitParam, sendParams model.SnmpSendParams) {

	queueInit, errOpen := amqp_serv.RabbitQueueInit(param.AmqpUrl, param.QueueName)
	if errOpen != nil {
		model.FailOnError(errOpen, "RabbitQueueOpen - FATAL ERROR")
	}

	defer queueInit.Connect.Close()
	defer queueInit.Channel.Close()

	// fmt.Println(queueInit)

	switch param.AmqpFuncType {
		case "get" :
			BasisGetAction(queueInit, sendParams, param)

		case "consumer" :
			BasisConsumeAction(queueInit, sendParams, param)

	    default:
			BasisGetAction(queueInit, sendParams, param)
	}

}


func BasisGetAction(queueInit  amqp_serv.QueueInitResultParam,
	                sendParams model.SnmpSendParams,
	                param model.CommonInitParam) {

	for {
		msg, err, ok := amqp_serv.GetOneMessage(queueInit)
		if err == nil && ok {
			snmpItem := amqp_serv.MessageDataConvert(msg, sendParams)
			go SnmpMakeRequest(snmpItem, param)
		}
	}

}


func BasisConsumeAction(queueInit  amqp_serv.QueueInitResultParam,
	                    sendParams model.SnmpSendParams,
	                    param model.CommonInitParam) {

	queueName := queueInit.Name

	messages, err := queueInit.Channel.Consume(
		queueName, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	if err != nil {
		model.WarnOnError(err, "Warn Error - Get Messages Consume Func")
	} else {

		ch := 0
		for msg := range messages {

			snmpItem := amqp_serv.MessageDataConvert(msg, sendParams)
			SnmpMakeRequest(snmpItem, param)

			fmt.Println("Num:", ch, "AmqpMessageItem:", snmpItem)
			ch++
		}
	}
}


/////////////////////////////////
/////////////////////////////////
//    Common Snmp Services

func SnmpMakeRequest(params model.SnmpSendParams, commonParam model.CommonInitParam) error {

	response, err := snmp_serv.SnmpManagerStart(params, commonParam.SnmpFuncType)
	if err != nil {
		model.WarnOnError(err, "Error: SnmpMakeRequest")
		SnmpRetriesAction(params)
		return err
	}

	jsonLog , jsonSaveError := snmp_serv.MakeJsonMultiRequest(commonParam.SaveApiUrl, response.Items)
	if jsonSaveError != nil {
		model.WarnOnError(err, "MakeJsonMultiRequest::Json SEND Error:")
	}


	fmt.Println("======= Start ========")
	fmt.Println("SnmpManagerStart : OK", response)
	fmt.Println("MakeJsonResult : OK", jsonLog)
	fmt.Println("======= End  ========")

	return jsonSaveError

}


func SnmpRetriesAction(params model.SnmpSendParams) error {

	retry := params.Retry
	if retry > 5 {
		_ = snmp_serv.SnmpExceptionHandler(params)
	} else {
		params.Retry++
		msg := amqp_serv.MessageConvertTostring(params)
		amqp_serv.ProducerAddMessage(msg)
	}

	return nil

}



////////////////////////////////////////
// New Version (Gorutines and channal)


func GetTasksFromQueue(messChannal chan<- model.SnmpSendParams,
	                   param model.CommonInitParam,
	                   sendParams model.SnmpSendParams,
                       wg *sync.WaitGroup) {


	queueInit, errOpen := amqp_serv.RabbitQueueInit(param.AmqpUrl, param.QueueName)
	if errOpen != nil {
		model.FailOnError(errOpen, "RabbitQueueOpen - FATAL ERROR")
	}
	defer queueInit.Connect.Close()
	defer queueInit.Channel.Close()

	// fmt.Println(queueInit)

	switch param.AmqpFuncType {
	case "get" :

		msg, err, ok := amqp_serv.GetOneMessage(queueInit)
		if err == nil && ok {
			snmpItem := amqp_serv.MessageDataConvert(msg, sendParams)
			messChannal <- snmpItem

		} else {
			ErrorGetTask(err)
		}

		//case "consumer" :
		//	BasisConsumeAction(queueInit, sendParams, param)
		//
		//default:
		//	BasisGetAction(queueInit, sendParams, param)

	}

	defer wg.Done()

}


func SendRequestToDevice(messChannal chan model.SnmpSendParams,
	                     param model.CommonInitParam,
                         wg *sync.WaitGroup) {

    message := <- messChannal
	SnmpMakeRequest(message, param)
    close(messChannal)
	defer wg.Done()


}


func ErrorGetTask(e error) {
	fmt.Println(e)
}