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
}

// NewAmqpHandler constructs the handler
func NewAmqpHandler(url string, logger logging.Logger) *AmqpHandler {
	return &AmqpHandler{url, logger, nil, nil}
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
