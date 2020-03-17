package amqp_handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	model "Snmp-Golang/internal/models"
	snmp_handl "Snmp-Golang/internal/snmp_handler"
	mq "github.com/streadway/amqp"
)

func RecevieMessagesListFromQueue(amqpUrl string, queueName string, saveApiUrl string) {

	connect, err := mq.Dial(amqpUrl)
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer connect.Close()

	channel, err := connect.Channel()
	FailOnError(err, "Failed to open a channel")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	FailOnError(err, "Failed to declare a queue")

	messagesList, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	// messagesList, _status, err := channel.Get(queue.Name, true)

	FailOnError(err, "Failed to register a consumer")

	sendParams := model.SnmpSendParams{
		Ip:  "192.168.2.184",
		Oid: ".1.3.6.1",
		// Oid:".1.3.6.1.2.1.1.9.1",
		Community: "public",
		Port:      "161",
		DeviceId:  "Pr-777-999",
		SelCount:  0,
	}

	//fmt.Println(string(messagesList.Body))
	//fmt.Println(sendParams)
	//fmt.Println(_status)

	MessagesListRendering(messagesList, sendParams, saveApiUrl)

	//streamCount := 10
	//
	//forever := make(chan bool)
	//
	//for a := 0; a < streamCount; a++ {
	//	go MessagesListRendering(messagesList, sendParams, saveApiUrl)
	//}
	//
	//log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	//
	//<-forever

}

func MessagesListRendering(messages <-chan mq.Delivery, sendParams model.SnmpSendParams, saveApiUrl string) {

	ch := 0
	for msg := range messages {

		amqpMessage := GetFormAmqpItem(msg, sendParams) // Формируем данные для запроса

		snmp_handl.SnmpBulkRequestSend(amqpMessage, saveApiUrl) // Выполняем запрос

		fmt.Println("Num:", ch, "AmqpMessageItem:", amqpMessage)

		ch++
	}

}

func GetFormAmqpItem(msg mq.Delivery, sendParams model.SnmpSendParams) model.SnmpSendParams {

	item := string(msg.Body)
	message := strings.Split(item, " ")

	sendParams.Id = message[0] // MessageId
	//sendParams.Ip   = message[1] // Ip Address
	//sendParams.Oid  = message[2] // Oid
	sendParams.Port = message[3]

	return sendParams
}

func FailOnError(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
		model.LogSave(msg, err, "FatalErrors")
		log.Fatalf("%s: %s", msg, err)
	}
}

func WarnOnError(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
		model.LogSave(msg, err, "WarningErrors")
	}
}

//////////////////////////////////
//////////////////////////////////

func RecevieGetMessageOne(amqpUrl string, queueName string) {

	connect, err := mq.Dial(amqpUrl)
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer connect.Close()

	channel, err := connect.Channel()
	FailOnError(err, "Failed to open a channel")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	FailOnError(err, "Failed to declare a queue")

	resultMessage, _, err := channel.Get(queue.Name, true)

	FailOnError(err, "Failed to register a consumer")

	message := string(resultMessage.Body)

	fmt.Println(message)

}

func MakeJsonRequest(apiUrl string, messages snmp_handl.SnmpResultItems) {


	bytesRepresentation, err := json.Marshal(messages.Items)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(apiUrl, "application/json",
		                   bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	var _result map[string]interface{}

	sendError := json.NewDecoder(resp.Body).Decode(&_result)

	log.Println(sendError)
	log.Println(_result)

}

func SnmpResultsRender(res snmp_handl.SnmpResultItems) {

	saveApiUrl := model.SAVE_API_URL
	messages := snmp_handl.SnmpResultItems{}

	for _, item := range res.Items {

		//messages = append(messages,
		//	              model.ResponseMessage{
		//					Oid:".1.3.6.1.4.1.119.2.3.69.501.7.1.1.1.3.17",
		//					Ip: "192.168.10.12 dff",
		//					Value:"Тест 20456",
		//					DeviceId: "234",
		//			      })
		fmt.Println(item)

	}

	MakeJsonRequest(saveApiUrl, messages)

}

//////////////////////////////////////////////
//////////////////////////////////////////////
///    NEW MESSAGES RENDER    ////////////////
//////////////////////////////////////////////

func GetMessagesListStart(amqpUrl, queueName, saveApiUrl, amqpFuncType string) {

	channel, connect, errChannel := AmqpChannelInit(amqpUrl)
	defer connect.Close()

	if errChannel != nil {
		FailOnError(errChannel, "Fatal Error - Connect to Cannel Amqp")
	}
	defer channel.Close()

	queue, errQueue := QueueDeclareInit(channel, queueName)
	if errQueue != nil {
		FailOnError(errQueue, "Fatal Error - Queue declare Amqp")
	}

	sendParams := model.SnmpSendParams{
		Ip:  "192.168.2.184",
		Oid: ".1.3.6.1",
		// Oid:".1.3.6.1.2.1.1.9.1",
		Community: "public",
		Port:      "161",
		DeviceId:  "Pr-777-999",
		SelCount:  0,
	}

	if amqpFuncType == "consumer" {

		messagesList, errGetMessages := channel.Consume(
			queue.Name, // queue
			"",         // consumer
			true,       // auto-ack
			false,      // exclusive
			false,      // no-local
			false,      // no-wait
			nil,        // args
		)
		if errGetMessages != nil {
			FailOnError(errGetMessages, "Fatal Error - Get Messages ConsumeFunc")
		}
		QueueMessagesListRender(messagesList, sendParams, saveApiUrl)

	} else {
		ch := 0
		for {

			message, stateOk, err := channel.Get(queue.Name, true)
			if err != nil {
				WarnOnError(err, "Message GetFunc (channel.Get):")
			}
			if stateOk {

				ch++
				fmt.Println("1================")
				fmt.Println("Message Count:", ch)
				model.DatetimePrint()
				fmt.Println("1================")

				QueueMessageExec(message, sendParams, saveApiUrl)
				// go QueueMessageExec(message, sendParams, saveApiUrl)

			} else {
				fmt.Println("Message GetFunc Not Ok :", stateOk)
				WarnOnError(nil, "Message GetFunc Not Ok :")
			}
		}
	}
}

func QueueMessagesListRender(messagesList <-chan mq.Delivery,
	                         sendParams model.SnmpSendParams,
	                         saveApiUrl string) {
	ch := 0
	for msg := range messagesList {
		QueueMessageExec(msg, sendParams, saveApiUrl)
		fmt.Println("Num:", ch)
		ch++
	}
}

func QueueMessageExec(msg mq.Delivery, sendParams model.SnmpSendParams, saveApiUrl string) {

	snmpFuncType := model.SNMP_FUNC_TYPE
	queueMessage := GetFormAmqpItem(msg, sendParams) // Формируем данные для запроса

	fmt.Println("==============Snmp Query Start===============")
	jsonSaveError := snmp_handl.SnmpNewStart(queueMessage, saveApiUrl, snmpFuncType)
	fmt.Println("QueueMessageExec:: QueueMessage  : ", queueMessage)
	fmt.Println("QueueMessageExec:: JsonSaveError : ", jsonSaveError)
	fmt.Println("==============Snmp Query Finish===============")
	// DatetimePrint()
}

func AmqpChannelInit(amqpApiUrl string) (*mq.Channel, *mq.Connection, error) {

	connect, err := mq.Dial(amqpApiUrl)
	if(err != nil) {
		FailOnError(err, "Fatal Error - Connect to RabbitMQ")
	}

	channel, err := connect.Channel()
	if err != nil {
		return  channel, connect, err
	}

	return channel, connect, nil
}

func QueueDeclareInit(channel *mq.Channel, queueName string) (mq.Queue, error) {

	queue, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	return queue, err
}

func DatetimePrint() {
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02d__%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	fmt.Println(formatted)
}
