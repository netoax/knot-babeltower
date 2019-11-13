package network

import (
	"github.com/CESARBR/knot-babeltower/pkg/logging"
)

const (
	exchangeFogOut = "FogOut"
	registerOutKey = "device.registered"
)

// MsgPublisher handle messages received from a service
type MsgPublisher struct {
	logger logging.Logger
	amqp   *AmqpHandler
}

// Publisher is the interface with methods that the publisher should have
type Publisher interface {
	SendRegisterDevice([]byte) error
}

// NewMsgPublisher constructs the MsgPublisher
func NewMsgPublisher(logger logging.Logger, amqp *AmqpHandler) *MsgPublisher {
	return &MsgPublisher{logger, amqp}
}

// SendRegisterDevice sends a registered message
func (mp *MsgPublisher) SendRegisterDevice(msg []byte) error {
	mp.logger.Debug("Sending register message")
	// TODO: receive message
	return mp.amqp.PublishPersistentMessage(exchangeFogOut, registerOutKey, msg)
}
