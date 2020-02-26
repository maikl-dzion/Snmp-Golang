package consumer

import (
	"../../internal/handlers"
	"github.com/streadway/amqp"
	"log"
	"strings"
)



type QueueInitParams struct {
	Durable   bool
	Delete    bool
	Exclusive bool
	NoWait    bool
	QueueName string
}


type ResponseMessage struct {
	Oid      string
	Value    string
	Ip       string
	DeviceId string
}



//func RabbitMQInit(amqpApiUrl string) (*amqp.Channel, error) {
//
//	connect , err := amqp.Dial(amqpApiUrl)
//
//	defer connect.Close()
//
//	if err != nil {
//		return 	nil , err
//	}
//
//
//	channel, err := connect.Channel()
//
//	defer channel.Close()
//
//	return channel , err
//
//}
//
//
//func GetQueueDeclare(ch *amqp.Channel,
//	                 queueParams QueueInitParams) (amqp.Queue, error) {
//
//	q, err := ch.QueueDeclare(
//		queueParams.QueueName,    // name of the queue
//		queueParams.Durable,
//		queueParams.Delete,
//		queueParams.Exclusive,
//		queueParams.NoWait,     // noWait
//		nil,               // arguments
//	)
//
//	return q, err
//}



func AmqpFullInit(amqpApiUrl string, queueParams QueueInitParams) (amqp.Queue, error) {


	connect, err := amqp.Dial(amqpApiUrl)
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer connect.Close()

	channel, err := connect.Channel()
	FailOnError(err, "Failed to open a channel")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		queueParams.QueueName,  // name of the queue
		queueParams.Durable,    // durable
		queueParams.Delete,     // delete when unused
		queueParams.Exclusive,  // exclusive
		queueParams.NoWait,     // no-Wait
		nil,               // arguments
	)

	return queue, err

}



func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}



func AmqpInitialize(amqpApiUrl string) (*amqp.Channel, error) {


	connect, err := amqp.Dial(amqpApiUrl)
	FailOnError(err, "Failed to connect to RabbitMQ")
	//defer connect.Close()

	channel, err := connect.Channel()
	FailOnError(err, "Failed to open a channel")
	// defer channel.Close()

	return channel, err

}


func QueueDeclareInit(ch *amqp.Channel, queueParams QueueInitParams) (amqp.Queue, error) {

	queue, err := ch.QueueDeclare(
		queueParams.QueueName,  // name of the queue
		queueParams.Durable,    // durable
		queueParams.Delete,     // delete when unused
		queueParams.Exclusive,  // exclusive
		queueParams.NoWait,     // no-Wait
		nil,               // arguments
	)

	return queue, err
}


//func RespMessagesInit() *[]ResponseMessage {
//	var messages []ResponseMessage
//	return &messages
//}
//
//
//
//func (p *[]ResponseMessage) FormMessage(oid, ipAddress, value, deviceId) {
//	mess := ResponseMessage {
//		    Oid:oid, Ip: ipAddress, Value:value,DeviceId: deviceId
//	}
//	p = append(p, mess)
//}


func ReceiveQueueMessages(channel *amqp.Channel, queueName string)  {


	messages, err := channel.Consume(queueName, "", true, false, false, false, nil)

	if err != nil {
		log.Println("Error consuming messages -> ", err)
	}

	// results := []ResponseMessage{}

	_count := 0


	port := "161"
	oid  := "1.3.6.1.2.1.1.1.0";
	ipAddress := "192.168.2.184"



	for message := range messages {

		mess := string(message.Body)
		item := strings.Split(mess, " ")

		//ipAddress := item[1]
		//oid       := item[2]
		//port := "161"

		//resp := ResponseMessage{ Oid:oid, Ip: ipAddress, Value:""}
		//results = append(results, resp)

		handler.SnmpRequestInit(ipAddress, oid, port)
		log.Printf("Count: %d, ItemId: %s", _count, item[0])
		_count++

	}


}