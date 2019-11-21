package network

import (
	"encoding/json"

	"github.com/CESARBR/knot-babeltower/pkg/entities"
	"github.com/CESARBR/knot-babeltower/pkg/logging"
)

const (
	queueNameFogIn  = "FogIn"
	exchangeFogIn   = "FogIn"
	bindingKeyFogIn = "device.*"
)

// Interactor is the use case to be executed
type Interactor interface {
	Execute(id string, args ...interface{}) error
}

// MsgConsumer handle messages received from a service
type MsgConsumer struct {
	logger        logging.Logger
	amqp          *AmqpHandler
	registerThing Interactor
}

func (mc *MsgConsumer) onMsgReceived(msgChan chan InMsg) {
	var thing entities.Thing
	for {
		msg := <-msgChan
		mc.logger.Debug("Message received:", string(msg.Body))

		switch msg.RoutingKey {
		case "device.register":
			err := json.Unmarshal(msg.Body, &thing)
			if err != nil {
				mc.logger.Error(err)
				continue
			}

			authorizationHeader := msg.Headers["Authorization"]
			err = mc.registerThing.Execute(thing.ID, thing.Name, authorizationHeader)
			if err != nil {
				mc.logger.Error(err)
				continue
			}
		}
	}
}

// NewMsgConsumer constructs the MsgConsumer
func NewMsgConsumer(logger logging.Logger, amqp *AmqpHandler, registerThing Interactor) *MsgConsumer {
	return &MsgConsumer{logger, amqp, registerThing}
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
