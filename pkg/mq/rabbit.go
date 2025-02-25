package mq

import (
	"context"
	"errors"
	"time"

	"github.com/golang/glog"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	RabbitMQModeWork = iota + 1
	RabbitMQModePubSub
	RabbitMQModeRouting
	RabbitMQModeTopics
)

type RabbitMQ struct {
	conn      *amqp.Connection
	ch        *amqp.Channel
	QueueName string
	Exchange  string
	Key       string
	Url       string

	// 1. Work Mode
	// 2. Publish/Subscribe Mode
	Mode uint8
}

func newRabbitMQ(queueName, exchange, key, url string) *RabbitMQ {
	return &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Url:       url,
	}
}

// Work Mode
func NewRabbitMQWorkClient(queueName string, url string) (*RabbitMQ, error) {
	mq := newRabbitMQ(queueName, "", "", url)
	var err error
	mq.conn, err = amqp.Dial(mq.Url)
	if err != nil {
		glog.Errorln("[RabbitMQ] Failed to connect to RabbitMQ: ", err.Error())
		return nil, err
	}
	mq.ch, err = mq.conn.Channel()
	if err != nil {
		glog.Errorln("[RabbitMQ] Failed to open a channel: ", err.Error())
		return nil, err
	}
	mq.Mode = RabbitMQModeWork
	return mq, nil
}

func (r *RabbitMQ) PublishSimple(message []byte) error {
	if r.Mode != RabbitMQModeWork {
		return errors.New("rabbitmq: mode mismatch")
	}
	// Try to declare the exchange,
	// build the queue if it's not existed, else do nothing.
	q, err := r.ch.QueueDeclare(
		r.QueueName, // name
		true,        // durable (queue)
		false,       // auto-deleted
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		glog.Errorln("[RabbitMQ] Failed to declare a queue: ", err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = r.ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // Message durability
			ContentType:  "text/plain",
			Body:         message,
		},
	)
	return nil
}

// Publish/Subscribe Mode
func NewRabbitMQPubSubClient(exchangeName string, url string) (*RabbitMQ, error) {
	mq := newRabbitMQ("", exchangeName, "", url)
	var err error
	mq.conn, err = amqp.Dial(mq.Url)
	if err != nil {
		glog.Errorln("[RabbitMQ] Failed to connect to RabbitMQ: ", err.Error())
		return nil, err
	}
	mq.ch, err = mq.conn.Channel()
	if err != nil {
		glog.Errorln("[RabbitMQ] Failed to open a channel: ", err.Error())
		return nil, err
	}
	mq.Mode = RabbitMQModePubSub
	return mq, nil
}

func (r *RabbitMQ) PublishPubSub(message []byte) error {
	if r.Mode != RabbitMQModePubSub {
		return errors.New("rabbitmq: mode mismatch")
	}
	// Try to declare the exchange,
	// build the queue if it's not existed, else do nothing.
	err := r.ch.ExchangeDeclare(
		r.Exchange, // name
		"fanout",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		glog.Errorln("[RabbitMQ] Failed to declare an exchange: ", err.Error())
		return err
	}

	// Send the message to exchange.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = r.ch.PublishWithContext(
		ctx,
		r.Exchange, // exchange
		"",         // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		glog.Errorln("[RabbitMQ] Failed to publish a message: ", err.Error())
		return err
	}
	return nil
}

// func (r *RabbitMQ) Recieve() []byte {
// 	// Try to declare the exchange.
// 	err := r.ch.ExchangeDeclare(
// 		r.Exchange,	// name
// 		"fanout",	// type
// 		true,	// durable
// 		false,	// auto-deleted
// 		false,	// internal
// 		false,	// no-wait
// 		nil,	// arguments
// 	)
// 	if err != nil {
// 		glog.Errorln("[RabbitMQ] Failed to declare an exchange: ", err.Error())
// 		return
// 	}

// 	q, err := r.ch.QueueDeclare(
// 		"",	// name
// 		false,	// durable
// 		false,	// delete when unused
// 		true,	// exclusive
// 		false,	// no-wait
// 		nil,	// arguments
// 	)
// 	if err != nil {
// 		glog.Errorln("[RabbitMQ] Failed to declare an queue: ", err.Error())
// 		return
// 	}

// 	err = r.ch.QueueBind(
// 		q.Name,	// queue name
// 		"", 	// routing key
// 		r.Exchange,	// exchange
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		glog.Errorln("[RabbitMQ] Failed to bind a queue: ", err.Error())
// 		return nil
// 	}

// 	msgs, err := r.ch.Consume(
// 		q.Name,	// queue
// 		"",	// consumer
// 		true,	// auto-ack
// 		false,	// exclusive
// 		false,	// no-local
// 		false,	// no-wait
// 		nil,	// arguments
// 	)
// 	if err != nil {
// 		glog.Errorln("[RabbitMQ] Failed to register a consumer: ", err.Error())
// 		return
// 	}

// 	forever := make(chan bool)

// 	go func() {
// 		for d := range msgs {
// 			// TODO do something here.
// 		}
// 	}()
// 	<-forever
// }

func (r *RabbitMQ) Destory() {
	r.ch.Close()
	r.conn.Close()
}
