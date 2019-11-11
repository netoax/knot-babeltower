package network

import "github.com/CESARBR/knot-babeltower/pkg/logging"

// AmqpHandler handles the connection, queues and exchanges declared
type AmqpHandler struct {
	url    string
	logger logging.Logger
}

// NewAmqpHandler constructs the handler
func NewAmqpHandler(url string, logger logging.Logger) *AmqpHandler {
	return &AmqpHandler{url, logger}
}

// Start starts the handler
func (ah *AmqpHandler) Start() {
	ah.logger.Debug("AMQP handler started")
	// TODO: Start amqp connection
}
