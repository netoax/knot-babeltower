package network

import (
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

func (ah *AmqpHandler) notifyWhenClosed() {
	errReason := <-ah.conn.NotifyClose(make(chan *amqp.Error))
	ah.logger.Infof("AMQP connection closed: %s", errReason)
	// TODO: try to reconnect
}

// Start starts the handler
func (ah *AmqpHandler) Start(started chan bool) {
	conn, err := amqp.Dial(ah.url)
	if err != nil {
		// TODO: try to reconnect
		ah.logger.Error(err)
		started <- false
		return
	}

	ah.conn = conn
	go ah.notifyWhenClosed()

	channel, err := conn.Channel()
	if err != nil {
		// TODO: try to create channel again
		ah.logger.Error(err)
		started <- false
		return
	}

	ah.logger.Debug("AMQP handler connected")
	ah.channel = channel
	started <- true
}
