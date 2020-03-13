package amqp_handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	model "../models"
	snmp_handl "../snmp_handler"
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
		queueName,       // name
		false,    // durable
		false,  // delete when unused
		false,   // exclusive
		false,    // no-wait
		nil,        // arguments
	)

	FailOnError(err, "Failed to declare a queue")

	messagesList, err := channel.Consume(
		queue.Name,     // queue
		"",    // consumer
		true,   // auto-ack
		false, // exclusive
		false,  // no-local
		false,  // no-wait
		nil,      // args
	)

	// messagesList, _status, err := channel.Get(queue.Name, true)

	FailOnError(err, "Failed to register a consumer")


	sendParams := model.SnmpSendParams{
		Ip:"192.168.2.184",
		Oid:".1.3.6.1",
		// Oid:".1.3.6.1.2.1.1.9.1",
		Community:"public",
		Port:"161",
		DeviceId:"Pr-777-999",
		SelCount:0,
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

		amqpMessage := GetFormAmqpItem(msg, sendParams)  // Формируем данные для запроса

		snmp_handl.SnmpBulkRequestSend(amqpMessage, saveApiUrl) // Выполняем запрос

		fmt.Println("Num:", ch, "AmqpMessageItem:", amqpMessage)

		ch++
	}

}

func GetFormAmqpItem(msg mq.Delivery, sendParams model.SnmpSendParams) model.SnmpSendParams {

	item    := string(msg.Body)
	message := strings.Split(item, " ")

	sendParams.Id   = message[0] // MessageId
	//sendParams.Ip   = message[1] // Ip Address
	//sendParams.Oid  = message[2] // Oid
	sendParams.Port = message[3]

	return sendParams;
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
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
		queueName,  // name
		false,      // durable
		false,    // delete when unused
		false,     // exclusive
		false,      // no-wait
		nil,          // arguments
	)

	FailOnError(err, "Failed to declare a queue")

	resultMessage , _ , err := channel.Get(queue.Name, true)

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

func GetMessagesListStart(amqpUrl, queueName, saveApiUrl, selectType string) {

	channel, errChan := AmqpChannelInit(amqpUrl)

	if errChan != nil {
		fmt.Println("Amqp Init :", errChan)
		defer channel.Close()
	}

	queue, errQueue := QueueDeclareInit(channel, queueName)

	if errQueue != nil {
		fmt.Println("Queue declare :", errQueue)
	}

	sendParams := model.SnmpSendParams{
		Ip:"192.168.2.184",
		Oid:".1.3.6.1",
		// Oid:".1.3.6.1.2.1.1.9.1",
		Community:"public",
		Port:"161",
		DeviceId:"Pr-777-999",
		SelCount:0,
	}

	if selectType == "consumer" {

		messagesList, errCon := channel.Consume(
			queue.Name,     // queue
			"",    // consumer
			true,   // auto-ack
			false, // exclusive
			false,  // no-local
			false,  // no-wait
			nil,      // args
		)

		if errCon != nil {
			fmt.Println("Messages ConsumeFunc :", errCon)
		}

		QueueMessagesListRender(messagesList, sendParams, saveApiUrl)

	} else {

		for {

			message, stateOk, err := channel.Get(queue.Name, true)

			if err != nil {
				fmt.Println("Message GetFunc :", err)
			}

			if stateOk {
				QueueMessageExec(message, sendParams, saveApiUrl)
			} else {
				fmt.Println("Message GetFunc Not Ok :", stateOk)
			}

		}
	}


	// messagesList, _status, err := channel.Get(queue.Name, true)

}


func QueueMessagesListRender(messagesList <-chan mq.Delivery,
	                         sendParams   model.SnmpSendParams,
	                         saveApiUrl   string) {

	ch := 0
	for msg := range messagesList {

		QueueMessageExec(msg, sendParams, saveApiUrl)

		fmt.Println("Num:", ch)
		ch++

	}

}


func QueueMessageExec(msg mq.Delivery, sendParams model.SnmpSendParams, saveApiUrl string)  {

	queueMessage := GetFormAmqpItem(msg, sendParams)  // Формируем данные для запроса

	// snmp_handl.SnmpBulkRequestSend(queueMessage, saveApiUrl) // Выполняем snmp запрос

	fmt.Println("QueueMessage : ", queueMessage)

	DatetimePrint()

}


func AmqpChannelInit(amqpApiUrl string) (*mq.Channel, error) {

	connect, err := mq.Dial(amqpApiUrl)
	FailOnError(err, "Failed to connect to RabbitMQ")
	//defer connect.Close()

	channel, err := connect.Channel()
	FailOnError(err, "Failed to open a channel")
	// defer channel.Close()

	return channel, err

}

func QueueDeclareInit(channel *mq.Channel, queueName string) (mq.Queue, error) {

	queue, err := channel.QueueDeclare(
		queueName,       // name
		false,    // durable
		false,  // delete when unused
		false,   // exclusive
		false,    // no-wait
		nil,        // arguments
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