package amqp_handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	model "../models"
	mq "github.com/streadway/amqp"
	snmp_handl "../snmp_handler"
)

type AmqpSendItem struct{
	Oid string
	Ip string
	Id string
	Port string
}

type AmqpSendItems struct{
	Items []AmqpSendItem
}

func MessagesListRendering(messages <-chan mq.Delivery) {

	//port := "161"
	//oid  := "1.3.6.1.2.1.1.1.0";
	//ipAddress := "192.168.2.184"

	sendParams := model.SnmpSendParams{
		"192.168.2.184",
		".1.3.6.1.2.1.1.1.0",
		"public",
		"161",
		0,
	}


	// items := AmqpSendItems{}
	ch    := 0
	for msg := range messages {

		amqpItem := FormAmqpItem(msg)

		//sendParams.Ip  = amqpItem.Ip
		//sendParams.Oid = amqpItem.Oid

		// items.Items = append(items.Items, res)
		snmp_handl.BulkRequestRun(sendParams)


		fmt.Println("I:", ch, "Mess:", amqpItem)

		ch++
	}

}


func RecevieMessagesListFromQueue(amqpUrl string, queueName string) {

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

	messages, err := channel.Consume(
		queue.Name,     // queue
		"",    // consumer
		true,   // auto-ack
		false, // exclusive
		false,  // no-local
		false,  // no-wait
		nil,      // args
	)

	FailOnError(err, "Failed to register a consumer")

	MessagesListRendering(messages)

	//forever := make(chan bool)
	//
	//
	//for a := 0; a < 10; a++ {
	//	go MessagesRendering(messages)
	//}
	//
	//
	//log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	//
	//
	//<-forever

}


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


func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}


func FormAmqpItem(msg mq.Delivery) AmqpSendItem {

	item    := string(msg.Body)
	message := strings.Split(item, " ")

	sendItem := AmqpSendItem {
		Id  : message[0],
		Ip  : message[1],
		Oid : message[2],
		Port: message[3],
	}

	return sendItem;

}


func MakeJsonRequest(apiUrl string, messages model.ResponseJsonItems) {


	bytesRepresentation, err := json.Marshal(messages)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(apiUrl, "application/json",
		                   bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)

	log.Println(result)
	log.Println(result["data"])

}



func SnmpResultsRender(res snmp_handl.SnmpResultItems) {

	saveApiUrl := model.SAVE_API_URL
	messages := model.ResponseJsonItems{}

	for _, item := range res.Items {

		messages = append(messages,
			              model.ResponseMessage{
							Oid:".1.3.6.1.4.1.119.2.3.69.501.7.1.1.1.3.17",
							Ip: "192.168.10.12 dff",
							Value:"Тест 20456",
							DeviceId: "234",
					      })

	}

	MakeJsonRequest(saveApiUrl, messages)

}