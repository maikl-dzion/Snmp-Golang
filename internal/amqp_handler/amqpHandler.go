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
	// snmp_handl "../snmp_handler"
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

	//sendParams := model.SnmpSendParams{
	//	"190.169.1.5",
	//	".1.3.6.1.4.1.119.2.3.69.501.7.1.1",
	//	"public",
	//	"161",
	//	0,
	//}

	ch    := 0
	items := AmqpSendItems{}

	for msg := range messages {

		res := FormAmqpItem(msg)
		items.Items = append(items.Items, res)
		// snmp_handl.BulkRequestRun(sendParams)
		fmt.Println("I:", ch)
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
		queueName,  // name
		false,      // durable
		false,    // delete when unused
		false,     // exclusive
		false,      // no-wait
		nil,          // arguments
	)

	FailOnError(err, "Failed to declare a queue")

	messages, err := channel.Consume(
		queue.Name,    // queue
		"",    // consumer
		true,   // auto-ack
		false, // exclusive
		false,  // no-local
		false,  // no-wait
		nil,      // args
	)

	FailOnError(err, "Failed to register a consumer")


	MessagesListRendering(messages)

	// fmt.Println(messages)

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



func SendCurlExec(saveApiUrl string, messages []model.ResponseMessage, messageId string) {

	payloadBytes, err := json.Marshal(messages)

	if err != nil {
		// handle err
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", saveApiUrl + "?" + messageId, body)
	if err != nil {
		// handle err
	}

	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer b7d03a6947b217efb6f3ec3bd3504582")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		// handle err
	}

	defer resp.Body.Close()

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