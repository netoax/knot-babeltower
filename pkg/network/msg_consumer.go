package network

import "github.com/CESARBR/knot-babeltower/pkg/logging"

const (
	queueNameFogIn  = "FogIn"
	exchangeFogIn   = "FogIn"
	bindingKeyFogIn = "device.*"
)

// MsgConsumer handle messages received from a service
type MsgConsumer struct {
	logger       logging.Logger
	amqp         *AmqpHandler
	msgPublisher *MsgPublisher
}

func (mc *MsgConsumer) onMsgReceived(msgChan chan InMsg) {
	for {
		msg := <-msgChan
		mc.logger.Debug("Message received:", string(msg.Body))

		switch msg.RoutingKey {
		case "device.register":
			err := mc.msgPublisher.SendRegisterDevice(msg.Body)
			if err != nil {
				mc.logger.Error(err)
				continue
			}
		}
	}
}

// NewMsgConsumer constructs the MsgConsumer
func NewMsgConsumer(logger logging.Logger, amqp *AmqpHandler, msgPublisher *MsgPublisher) *MsgConsumer {
	return &MsgConsumer{logger, amqp, msgPublisher}
}

// Start starts to listen messages
func (mc *MsgConsumer) Start(started chan bool) {
	mc.logger.Debug("Msg consumer started")
	err := mc.amqp.DeclareQueue(queueNameFogIn, exchangeFogIn)
	if err != nil {
		mc.logger.Error(err)
		started <- false
		return
	}

	msgChan := make(chan InMsg)
	err = mc.amqp.OnMessage(msgChan, queueNameFogIn, exchangeFogIn, bindingKeyFogIn)
	if err != nil {
		mc.logger.Error(err)
		started <- false
		return
	}

	go mc.onMsgReceived(msgChan)

	started <- true
}

// Stop stops to listen for messages
func (mc *MsgConsumer) Stop() {
	mc.logger.Debug("Msg consumer stopped")
}
