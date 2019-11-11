package network

import "github.com/CESARBR/knot-babeltower/pkg/logging"

// AmqpHandler handles the connection, queues and exchanges declared
type AmqpHandler struct {
	logger logging.Logger
}

// NewAmqpHandler constructs the handler
func NewAmqpHandler(logger logging.Logger) *AmqpHandler {
	return &AmqpHandler{logger}
}

// Start starts the handler
func (ah *AmqpHandler) Start() {
	ah.logger.Debug("AMQP handler started")
	// TODO: Start amqp connection
}
