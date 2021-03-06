package simpleamqp

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

type AMQPConsumer interface {
	Receive(exchange string, routingKeys []string, queue string, queueOptions QueueOptions, queueTimeout time.Duration) chan AmqpMessage
	ReceiveWithoutTimeout(exchange string, routingKeys []string, queue string, queueOptions QueueOptions) chan AmqpMessage
}

type AmqpConsumer struct {
	brokerURI string
}

// Return AMQP Consumer
func NewAmqpConsumer(brokerURI string) *AmqpConsumer {
	return &AmqpConsumer{
		brokerURI: brokerURI,
	}
}

// AmqpMessage struct
type AmqpMessage struct {
	Exchange   string
	RoutingKey string
	Body       string
}

// Receive Return a AmqpMessage channel to receive messages using a given queue connected to the exchange with one ore more routing keys
// Autoreconnect on error or when we have no message after queueTimeout expired. Use 0 when not timeout is required.
// The function declares the queue
func (client *AmqpConsumer) Receive(exchange string, routingKeys []string, queue string, queueOptions QueueOptions, queueTimeout time.Duration) chan AmqpMessage {
	output := make(chan AmqpMessage)

	conn, ch, qname := client.setupConsuming(exchange, routingKeys, queue, queueOptions)

	go func() {
		for {
			messages, err := ch.Consume(qname, "", true, false, false, false, nil)
			if err != nil {
				log.Println("[simpleamqp] Error consuming messages -> ", err)
			}

			for closed := false; closed != true; {
				closed = messageToOuput(messages, output, queueTimeout)
			}

			log.Println("[simpleamqp] Closing connection ...")
			ch.Close()
			conn.Close()

			log.Println("[simpleamqp] Waiting before reconnect")
			time.Sleep(timeToReconnect)

			conn, ch, qname = client.setupConsuming(exchange, routingKeys, queue, queueOptions)
		}
	}()

	return output
}

// ReceiveWithoutTimeout the same behavior that Receive method, but without using a timeout for receiving from the queue
func (client *AmqpConsumer) ReceiveWithoutTimeout(exchange string, routingKeys []string, queue string, queueOptions QueueOptions) chan AmqpMessage {
	return client.Receive(exchange, routingKeys, queue, queueOptions, 0*time.Second)
}

func (client *AmqpConsumer) setupConsuming(exchange string, routingKeys []string, queue string, queueOptions QueueOptions) (*amqp.Connection, *amqp.Channel, string) {
	conn, ch := Setup(client.brokerURI)

	exchangeDeclare(ch, exchange)

	q := QueueDeclare(ch, queue, queueOptions)

	for _, routingKey := range routingKeys {
		ch.QueueBind(q.Name, routingKey, exchange, false, nil)
	}
	return conn, ch, q.Name
}

func messageToOuput(messages <-chan amqp.Delivery, output chan AmqpMessage, queueTimeout time.Duration) (closed bool) {

	if queueTimeout == 0*time.Second {
		message, more := <-messages
		if more {
			output <- AmqpMessage{Exchange: message.Exchange, RoutingKey: message.RoutingKey, Body: string(message.Body)}
			return false
		}
		log.Println("[simpleamqp] No more messages... closing channel to reconnect with timeout zero")
		return true
	}

	timeoutTimer := time.NewTimer(queueTimeout)
	defer timeoutTimer.Stop()
	afterTimeout := timeoutTimer.C

	select {
	case message, more := <-messages:
		if more {
			output <- AmqpMessage{Exchange: message.Exchange, RoutingKey: message.RoutingKey, Body: string(message.Body)}
			return false
		}
		log.Println("[simpleamqp] No more messages... closing channel to reconnect")
		return true
	case <-afterTimeout:
		log.Println("[simpleamqp] Too much time without messages... closing channel to reconnect")
		return true
	}

}