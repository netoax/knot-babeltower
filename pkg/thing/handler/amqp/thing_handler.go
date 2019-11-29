package amqp

import (
	"encoding/json"
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/CESARBR/knot-babeltower/pkg/thing/interactors"
)

const (
	queueNameFogIn  = "FogIn"
	exchangeFogIn   = "FogIn"
	bindingKeyFogIn = "device.*"
)

// MsgConsumer handle messages received from a service
type MsgConsumer struct {
	logger          logging.Logger
	amqp            *network.AmqpHandler
	thingInteractor interactors.Interactor
}

// UpdateSchemaRequest represents the update schema msg
type UpdateSchemaRequest struct {
	ID     string            `json:"id"`
	Schema []entities.Schema `json:"schema,omitempty"`
}

// NewMsgConsumer constructs the MsgConsumer
func NewMsgConsumer(logger logging.Logger, amqp *network.AmqpHandler, thingInteractor interactors.Interactor) *MsgConsumer {
	return &MsgConsumer{logger, amqp, thingInteractor}
}

func (mc *MsgConsumer) onMsgReceived(msgChan chan network.InMsg) {
	var thing entities.Thing
	var updateSchemaReq UpdateSchemaRequest
	for {
		msg := <-msgChan
		mc.logger.Debug("Message received:", string(msg.Body))

		switch msg.RoutingKey {
		case "device.register":
			err := json.Unmarshal(msg.Body, &thing)
			fmt.Print(thing)
			if err != nil {
				mc.logger.Error(err)
				continue
			}

			authorizationHeader := msg.Headers["Authorization"]
			err = mc.thingInteractor.Register(authorizationHeader.(string), thing.ID, thing.Name)
			if err != nil {
				mc.logger.Error(err)
				continue
			}
		case "schema.update":
			err := json.Unmarshal(msg.Body, &updateSchemaReq)
			if err != nil {
				mc.logger.Error(err)
				continue
			}

			authorizationHeader := msg.Headers["Authorization"]
			err = mc.thingInteractor.UpdateSchema(authorizationHeader.(string), updateSchemaReq.ID, updateSchemaReq.Schema)
			if err != nil {
				mc.logger.Error(err)
				continue
			}
		}
	}
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

	msgChan := make(chan network.InMsg)
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
