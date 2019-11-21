package interactors

import (
	"encoding/json"
	"strconv"

	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
)

// RegisterThing use case to register a new thing
type RegisterThing struct {
	logger       logging.Logger
	msgPublisher network.Publisher
	thingProxy   network.ThingProxy
}

type registerResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	Error error  `json:"error"`
}

// ErrorIDLenght is raised when ID is more than 16 characters
type ErrorIDLenght struct{}

// ErrorIDInvalid is raised when ID is not in hexadecimal value
type ErrorIDInvalid struct{}

// ErrorNameNotFound is raised when Name is empty
type ErrorNameNotFound struct{}

// ErrorArgument is raised when Name is empty
type ErrorArgument struct{ msg string }

func (err ErrorIDLenght) Error() string {
	return "ID lenght error"
}

func (err ErrorIDInvalid) Error() string {
	return "ID is not in hexadecimal"
}

func (err ErrorNameNotFound) Error() string {
	return "Name not found"
}

func (err ErrorArgument) Error() string {
	return err.msg
}

// NewRegisterThing contructs the use case
func NewRegisterThing(logger logging.Logger, msgPublisher network.Publisher, thingProxy network.ThingProxy) *RegisterThing {
	return &RegisterThing{logger, msgPublisher, thingProxy}
}

func (rt *RegisterThing) verifyThingID(id string) error {
	if len(id) > 16 {
		return ErrorIDLenght{}
	}

	_, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		rt.logger.Error(err)
		return ErrorIDInvalid{}
	}

	return nil
}

func (rt *RegisterThing) verifyArguments(args ...interface{}) error {
	if len(args) < 1 {
		return ErrorArgument{"Missing argument name"}
	}

	name, ok := args[0].(string)
	if !ok {
		return ErrorArgument{msg: "Name is not string"}
	}

	if len(name) == 0 {
		return ErrorNameNotFound{}
	}

	return nil
}

// Execute runs the use case
func (rt *RegisterThing) Execute(id string, args ...interface{}) error {
	rt.logger.Debug("Executing register thing use case")
	err := rt.verifyArguments(args...)
	if err != nil {
		return err
	}

	err = rt.verifyThingID(id)

	// TODO: add proxy request to token
	response := registerResponse{ID: id, Token: "secret", Error: err}
	bytes, err := json.Marshal(response)
	if err != nil {
		rt.logger.Error(err)
		return err
	}

	err = rt.msgPublisher.SendRegisterDevice(bytes)
	if err != nil {
		rt.logger.Error(err)
		return err
	}

	return nil
}
