package network

import (
	"gopkg.in/cenkalti/backoff.v3"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/streadway/amqp"
)

// AmqpHandler handles the connection, queues and exchanges declared
type AmqpHandler struct {
	url     string
	logger  logging.Logger
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
}

// NewAmqpHandler constructs the handler
func NewAmqpHandler(url string, logger logging.Logger) *AmqpHandler {
	return &AmqpHandler{url, logger, nil, nil, nil}
}

func (ah *AmqpHandler) notifyWhenClosed(started chan bool) {
	errReason := <-ah.conn.NotifyClose(make(chan *amqp.Error))
	ah.logger.Infof("AMQP connection closed: %s", errReason)
	started <- false
	if errReason != nil {
		err := backoff.Retry(ah.connect, backoff.NewExponentialBackOff())
		if err != nil {
			ah.logger.Error(err)
			started <- false
			return
		}

		go ah.notifyWhenClosed(started)
		started <- true
	}
}

func (ah *AmqpHandler) connect() error {
	conn, err := amqp.Dial(ah.url)
	if err != nil {
		ah.logger.Error(err)
		return err
	}

	ah.conn = conn

	channel, err := ah.conn.Channel()
	if err != nil {
		ah.logger.Error(err)
		return err
	}

	ah.logger.Debug("AMQP handler connected")
	ah.channel = channel

	return nil
}

// Start starts the handler
func (ah *AmqpHandler) Start(started chan bool) {
	err := backoff.Retry(ah.connect, backoff.NewExponentialBackOff())
	if err != nil {
		ah.logger.Error(err)
		started <- false
		return
	}

	go ah.notifyWhenClosed(started)
	started <- true
}

// DeclareQueue declare a queue in amqp handler
func (ah *AmqpHandler) DeclareQueue(queueName, exchangeName string) error {
	err := ah.channel.ExchangeDeclare(
		exchangeName,
		amqp.ExchangeTopic, // type
		true,               // durable
		false,              // delete when complete
		false,              // internal
		false,              // noWait
		nil,                // arguments
	)
	if err != nil {
		ah.logger.Error(err)
		return err
	}

	queue, err := ah.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		ah.logger.Error(err)
		return err
	}

	ah.queue = &queue
	return nil
}

// OnMessage receive messages an put them on channel
func (ah *AmqpHandler) OnMessage(msgChan chan InMsg, queueName, exchange, key string) error {

	err := ah.channel.QueueBind(
		queueName,
		key,
		exchange,
		false,
		nil,
	)
	if err != nil {
		ah.logger.Error(err)
		return err
	}

	deliveries, err := ah.channel.Consume(
		queueName,
		"consumer-tag",
		true,  // noAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		ah.logger.Error(err)
		return err
	}

	go func(deliveries <-chan amqp.Delivery) {
		for d := range deliveries {
			msgChan <- InMsg{d.Exchange, d.RoutingKey, d.Headers, d.Body}
		}
	}(deliveries)

	return nil
}

// PublishPersistentMessage sends a persistent message to RabbitMQ
func (ah *AmqpHandler) PublishPersistentMessage(exchange, key string, body []byte) error {
	err := ah.channel.ExchangeDeclare(
		exchange,
		amqp.ExchangeTopic, // type
		true,               // durable
		false,              // delete when complete
		false,              // internal
		false,              // noWait
		nil,                // arguments
	)
	if err != nil {
		ah.logger.Error(err)
		return err
	}

	err = ah.channel.Publish(
		exchange,
		key,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            body,
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
		},
	)
	if err != nil {
		ah.logger.Error(err)
		return err
	}

	return nil
}

// Stop closes the connection started
func (ah *AmqpHandler) Stop() {
	if ah.conn != nil && !ah.conn.IsClosed() {
		ah.conn.Close()
	}

	if ah.channel != nil {
		ah.channel.Close()
	}

	ah.logger.Debug("AMQP handler stopped")
}
