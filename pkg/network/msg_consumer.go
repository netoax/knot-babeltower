package network

import "github.com/CESARBR/knot-babeltower/pkg/logging"

// MsgConsumer handle messages received from a service
type MsgConsumer struct {
	logger logging.Logger
	amqp   *AmqpHandler
}

// NewMsgConsumer constructs the MsgConsumer
func NewMsgConsumer(logger logging.Logger, amqp *AmqpHandler) *MsgConsumer {
	return &MsgConsumer{logger, amqp}
}

// Start starts to listen messages
func (mc *MsgConsumer) Start(started chan bool) {
	mc.logger.Debug("Msg consumer started")
	started <- true
}

// Stop stops to listen for messages
func (mc *MsgConsumer) Stop() {
	mc.logger.Debug("Msg consumer stopped")
}
