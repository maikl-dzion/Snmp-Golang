package amqp_services

import (
	model "Snmp-Golang/internal/models"
	"github.com/streadway/amqp"
	"strconv"
	"strings"
)

type QueueInitResultParam struct {
	Channel  *amqp.Channel
	Connect  *amqp.Connection
	Queue    *amqp.Queue
	Name     string
}


func AmqpChannelInit(amqpApiUrl string) (QueueInitResultParam, error) {

	connect, err := amqp.Dial(amqpApiUrl)
	if err != nil {
		model.FailOnError(err, "Fatal Error - Connect to RabbitMQ")
	}

	queueInit := QueueInitResultParam{}

	channel, err := connect.Channel()
	if err != nil {
		return queueInit, err
	}

	queueInit.Connect = connect
	queueInit.Channel = channel

	return queueInit, nil
}

func QueueDeclareInit(param *QueueInitResultParam, queueName string) error {

	queue, err := param.Channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		model.WarnOnError(err, "Warn Error - Not QueueDeclare Amqp")
		return err
	}

	param.Queue = &queue
	param.Name  = queue.Name

	return nil
}

func RabbitQueueInit(amqpUrl, queueName string) (QueueInitResultParam, error) {

	queueInit , errChannel := AmqpChannelInit(amqpUrl)

	if errChannel != nil {
		model.WarnOnError(errChannel, "Fatal Error - Connect to Cannel Amqp")
		return queueInit, errChannel
	}

	errQueue := QueueDeclareInit(&queueInit, queueName)
	if errQueue != nil {
		model.WarnOnError(errQueue, "Fatal Error - Queue declare Amqp")
		return queueInit, errQueue
	}

	return queueInit, nil
}




func GetOneMessage(init QueueInitResultParam) (amqp.Delivery, error, bool) {

	queueName := init.Name
	msg, statOk , err := init.Channel.Get(queueName, true)
	if err != nil {
		model.WarnOnError(err, "Not Get message in Queues - WARN ERROR")
		return msg, err, false
	}

	if statOk {
		return msg, nil, true
	} else {
		return msg, nil, false
	}

}



func MessageDataConvert(msg amqp.Delivery,
	                    sendParams model.SnmpSendParams) model.SnmpSendParams {

	item := string(msg.Body)
	message := strings.Split(item, " ")

	sendParams.Id = message[0]                     // MessageId
	sendParams.Ip = message[1]                     // Ip Address
	sendParams.Oid = message[2]                    // Oid
	sendParams.Port = message[3]                   // Port
	sendParams.Type = message[4]                   // Type
	sendParams.Retry, _ = strconv.Atoi(message[5]) // Retry

	return sendParams
}


func MessageConvertTostring(sendParams model.SnmpSendParams) string {

	retry := string(sendParams.Retry)

	msg := sendParams.Id  + " " +
		   sendParams.Ip  + " " +
		   sendParams.Oid + " " +
		   sendParams.Port + " " +
		   sendParams.Type + " " +
		   retry

	return msg
}


func ProducerAddMessage(msg string) error {

	config    := model.GetCommonInitParam()
	amqpUrl   := config.AmqpUrl
	queueName := config.QueueName

	conn, _ := amqp.Dial(amqpUrl)
	defer conn.Close()

	//---Create a channel
	ch, _ := conn.Channel()
	defer ch.Close()

	//---Declare a queue
	q, err := ch.QueueDeclare(
		queueName, // name of the queue
		false,     // should the message be persistent? also queue will survive if the cluster gets reset
		false,     // autodelete if there's no consumers (like queues that have anonymous names, often used with fanout exchange)
		false,     // exclusive means I should get an error if any other consumer subsribes to this queue
		false,     // no-wait means I don't want RabbitMQ to wait if there's a queue successfully setup
		nil,       // arguments for more advanced configuration
	)

	if err != nil {
		model.WarnOnError(err, "AmqpProducer::QueueDeclare ERROR:")
		return err
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)


	eventName   := "Retry"
	message     := "Retry Message In RabbitMq"
	logFileName := "retries_in_queue"
	model.EventLogSave(message, eventName, logFileName)

	return err
}



