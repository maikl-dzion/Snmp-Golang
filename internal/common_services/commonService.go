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

	//r := snmp_serv.SnmpResultItems{}
	//for i := 0; i < 15; i++ {
	//
	//	item := snmp_serv.SnmpResultMessage{}
	//	item.Oid = "0.898.46756"
	//	item.DataType = "0.898.46756"
	//	item.ValueStr = "0.898.46756"
	//	item.ValueInt = 567
	//	item.DeviceId = "FFFGGGG"
	//	item.Ip = "192.45.56.57"
	//
	//	r.Items = append(r.Items, item)
	//}
	//response.Items = r.Items

	jsonLog , jsonSaveError := snmp_serv.MakeJsonMultiRequest(commonParam.SaveApiUrl, response.Items)
	if jsonSaveError != nil {
		model.WarnOnError(err, "MakeJsonMultiRequest::Json SEND Error:")
	}

	printSnmpActionResultView(response, jsonLog)

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



func printSnmpActionResultView(rSnmp snmp_serv.SnmpResultItems,
	                           rJson map[string]interface{}) {
	//fmt.Println("======= Start ========")
	//fmt.Println("SnmpManagerStart : OK", rSnmp)
	//fmt.Println("MakeJsonResult : OK"  , rJson)
	//fmt.Println("======= End  ========")
}

////////////////////////////////////////
// New Version (Gorutines and channal)


func GetTasksFromQueue(messChannal chan<- model.SnmpSendParams,
	                   param model.CommonInitParam,
	                   sendParams model.SnmpSendParams,
                       wg *sync.WaitGroup, taskCount int) {


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
			printTaskCount(taskCount, 1)

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
                         wg *sync.WaitGroup, taskCount int) {

    message := <- messChannal
	SnmpMakeRequest(message, param)
	printTaskCount(taskCount, 2)
	// time.Sleep(1 * time.Second)
    close(messChannal)
	defer wg.Done()

}


func ErrorGetTask(e error) {
	fmt.Println(e)
}


func printTaskCount(ch int, stype int) {

	switch stype {

	case 1 :
		fmt.Println("Кол. полученных заданий", ch)

	case 2 :
		fmt.Println("Кол. выполненных заданий", ch)

	}
}







///////////////////////////////
///////////////////////////////


func GetTasksFromQueueNew(param model.CommonInitParam,
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
			SendRequestToDeviceNew(snmpItem, param)

		} else {
			ErrorGetTask(err)
		}

	}

	defer wg.Done()

}


func SendRequestToDeviceNew(message model.SnmpSendParams,
	                        param   model.CommonInitParam) {

	SnmpMakeRequestNew(message, param)
	// printTaskCount(taskCount, 2)
	// close(messChannal)
	// defer wg.Done()

}


func SnmpMakeRequestNew(params model.SnmpSendParams,
	                    commonParam model.CommonInitParam) error {

	response, err := snmp_serv.SnmpManagerStart(params, commonParam.SnmpFuncType)
	if err != nil {
		model.WarnOnError(err, "Error: SnmpMakeRequest")
		SnmpRetriesAction(params)
		return err
	}

	//r := snmp_serv.SnmpResultItems{}
	//for i := 0; i < 15; i++ {
	//
	//	item := snmp_serv.SnmpResultMessage{}
	//	item.Oid = "0.898.46756"
	//	item.DataType = "0.898.46756"
	//	item.ValueStr = "0.898.46756"
	//	item.ValueInt = 567
	//	item.DeviceId = "FFFGGGG"
	//	item.Ip = "192.45.56.57"
	//
	//	r.Items = append(r.Items, item)
	//}
	//response.Items = r.Items


	jsonLog , jsonSaveError := snmp_serv.MakeJsonMultiRequest(commonParam.SaveApiUrl, response.Items)
	if jsonSaveError != nil {
		model.WarnOnError(err, "MakeJsonMultiRequest::Json SEND Error:")
	}

	printSnmpActionResultView(response, jsonLog)

	return jsonSaveError

}
