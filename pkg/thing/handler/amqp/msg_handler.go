package amqp

import (
	"encoding/json"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/CESARBR/knot-babeltower/pkg/thing/interactors"
)

const (
	queueNameFogIn   = "fogIn-messages"
	exchangeFogIn    = "fogIn"
	bindingKeyDevice = "device.*"
	bindingKeySchema = "schema.*"
)

// MsgHandler handle messages received from a service
type MsgHandler struct {
	logger          logging.Logger
	amqp            *network.Amqp
	thingInteractor interactors.Interactor
}

// UpdateSchemaRequest represents the update schema msg
type UpdateSchemaRequest struct {
	ID     string            `json:"id"`
	Schema []entities.Schema `json:"schema,omitempty"`
}

func (mc *MsgHandler) handleRegisterMsg(body []byte, authorizationHeader string) error {
	msgParsed := network.RegisterRequestMsg{}
	err := json.Unmarshal(body, &msgParsed)
	if err != nil {
		return err
	}

	return mc.thingInteractor.Register(msgParsed.ID, msgParsed.Name, authorizationHeader)
}

func (mc *MsgHandler) onMsgReceived(msgChan chan network.InMsg) {
	var updateSchemaReq UpdateSchemaRequest
	for {
		msg := <-msgChan
		mc.logger.Infof("Exchange: %s, routing key: %s", msg.Exchange, msg.RoutingKey)
		mc.logger.Infof("Message received: %s", string(msg.Body))

		authorizationHeader := msg.Headers["Authorization"]

		switch msg.RoutingKey {
		case "device.register":
			err := mc.handleRegisterMsg(msg.Body, authorizationHeader.(string))
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
			mc.logger.Info("Update schema message received")
			mc.logger.Info(authorizationHeader, updateSchemaReq)
		}
	}
}

// NewMsgHandler constructs the MsgHandler
func NewMsgHandler(logger logging.Logger, amqp *network.Amqp, registerThing interactors.Interactor) *MsgHandler {
	return &MsgHandler{logger, amqp, registerThing}
}

func (mc *MsgHandler) subscribeToMessages(msgChan chan network.InMsg) error {
	var err error
	subscribe := func(msgChan chan network.InMsg, queueName, exchange, key string) {
		if err != nil {
			return
		}
		err = mc.amqp.OnMessage(msgChan, queueName, exchange, key)
	}

	subscribe(msgChan, queueNameFogIn, exchangeFogIn, bindingKeyDevice)
	subscribe(msgChan, queueNameFogIn, exchangeFogIn, bindingKeySchema)
	return err
}

// Start starts to listen messages
func (mc *MsgHandler) Start(started chan bool) {
	mc.logger.Debug("Msg handler started")
	msgChan := make(chan network.InMsg)
	err := mc.subscribeToMessages(msgChan)
	if err != nil {
		mc.logger.Error(err)
		started <- false
		return
	}

	go mc.onMsgReceived(msgChan)

	started <- true
}

// Stop stops to listen for messages
func (mc *MsgHandler) Stop() {
	mc.logger.Debug("Msg handler stopped")
}
