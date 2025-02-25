package main

import (
	"context"
	checkoutpb "goshop/api/protobuf/checkout"

	"github.com/golang/glog"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

func rabbitConsumer(queueName string, url string) {
	conn, err := amqp.Dial(url)
	if err != nil {
		glog.Errorln("[CheckoutServer] Dial RabbitMQ error: ", err.Error())
		return
	}

	ch, err := conn.Channel()
	if err != nil {
		glog.Errorln("[CheckoutServer] RabbitMQ error: ", err.Error())
		return
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		glog.Errorln("[CheckoutServer] RabbitMQ declare error: ", err.Error())
		return
	}

	// Fair dispatch
	err = ch.Qos(
		1,     // prefetch count, most message given to a worker at a time.
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		glog.Errorln("[CheckoutServer] RabbitMQ failed to set Qos: ", err.Error())
		return
	}

	messages, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		glog.Errorln("[CheckoutServer] RabbitMQ failed to register a consumer: ", err.Error())
		return
	}

	var forever chan struct{}

	checkoutService := new(CheckoutRpcService)
	req := &checkoutpb.ReqCheckout{}
	go func() {
		for d := range messages {
			err = proto.Unmarshal(d.Body, req)
			if err != nil {
				glog.Errorln("[CheckoutServer] RabbitMQ failed to unmarshal the message: ", err.Error())
				continue
			}
			// FIXME: dial service timeout
			checkoutService.Checkout(context.Background(), req)

			// NOTE: Message acknowledgment(ACK).
			// If a consumer dies (its channel is closed, connection is closed, or TCP connection is lost) without sending an ack,
			// RabbitMQ will understand that a message wasn't processed fully and will re-queue it.
			// If there are other consumers online at the same time, it will then quickly redeliver it to another consumer.
			// That way you can be sure that no message is lost, even if the workers occasionally die.
			d.Ack(false)
		}
	}()

	<-forever
}
