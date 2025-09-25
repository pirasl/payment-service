package main

import amqp "github.com/rabbitmq/amqp091-go"

type rabbitMQClient struct {
	ampqpUrl string
	conn     *amqp.Connection
	channel  *amqp.Channel
}

func newRabbitMQClient() (*rabbitMQClient, error) {
	rabbitMQUrl, err := getRequiredStringEnv("AMQP_URL")
	if err != nil {
		return nil, err
	}

	conn, err := amqp.Dial(*rabbitMQUrl)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"stripe_events", // name
		"fanout",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"stripe_processing",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,
		"",
		"stripe_events",
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &rabbitMQClient{
		ampqpUrl: *rabbitMQUrl,
		conn:     conn,
		channel:  ch,
	}, nil

}
