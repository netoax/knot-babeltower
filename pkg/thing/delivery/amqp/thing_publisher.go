package amqp

import (
	"encoding/json"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
)

const (
	exchangeFogOut = "FogOut"
	registerOutKey = "device.registered"
	schemaOutKey   = "schema.updated"
)

// MsgPublisher handle messages received from a service
type MsgPublisher struct {
	logger logging.Logger
	amqp   *network.AmqpHandler
}

// UpdatedSchemaResponse represents the update schema response mapped from use case to the AMQP response
type UpdatedSchemaResponse struct {
	ID string `json:"id"`
}

// Publisher is the interface with methods that the publisher should have
type Publisher interface {
	SendRegisterDevice([]byte) error
	SendUpdatedSchema(thingID string) error
}

// NewMsgPublisher constructs the MsgPublisher
func NewMsgPublisher(logger logging.Logger, amqp *network.AmqpHandler) *MsgPublisher {
	return &MsgPublisher{logger, amqp}
}

// SendRegisterDevice sends a registered message
func (mp *MsgPublisher) SendRegisterDevice(msg []byte) error {
	mp.logger.Debug("Sending register message")
	// TODO: receive message
	return mp.amqp.PublishPersistentMessage(exchangeFogOut, registerOutKey, msg)
}

// SendUpdatedSchema sends the updated schema response
func (mp *MsgPublisher) SendUpdatedSchema(thingID string) error {
	resp := &UpdatedSchemaResponse{thingID}
	msg, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	return mp.amqp.PublishPersistentMessage(exchangeFogOut, schemaOutKey, msg)
}
