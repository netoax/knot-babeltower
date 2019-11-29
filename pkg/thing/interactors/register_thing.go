package interactors

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type registerResponse struct {
	ID    string  `json:"id"`
	Token *string `json:"token"`
	Error *string `json:"error"`
}

// ErrorIDLenght is raised when ID is more than 16 characters
type ErrorIDLenght struct{}

// ErrorIDInvalid is raised when ID is not in hexadecimal value
type ErrorIDInvalid struct{}

// ErrorNameNotFound is raised when Name is empty
type ErrorNameNotFound struct{}

// ErrorUnauthorized is raised when authorization token is empty
type ErrorUnauthorized struct{}

// ErrorMissingArgument is raised there is some argument missing
type ErrorMissingArgument struct{}

// ErrorInvalidTypeArgument is raised when the type is the expected
type ErrorInvalidTypeArgument struct{ msg string }

func (err ErrorIDLenght) Error() string {
	return "ID lenght error"
}

func (err ErrorIDInvalid) Error() string {
	return "ID is not in hexadecimal"
}

func (err ErrorNameNotFound) Error() string {
	return "Name not found"
}

func (err ErrorUnauthorized) Error() string {
	return "Authorization token not found"
}

func (err ErrorMissingArgument) Error() string {
	return "Missing arguments"
}

func (err ErrorInvalidTypeArgument) Error() string {
	return err.msg
}

func (i *ThingInteractor) verifyThingID(id string) error {
	if len(id) > 16 {
		return ErrorIDLenght{}
	}

	_, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		i.logger.Error(err)
		return ErrorIDInvalid{}
	}

	return nil
}

// Register runs the use case
func (i *ThingInteractor) Register(authorization, id, name string) error {
	token := ""
	i.logger.Debug("Executing register thing use case")
	err := i.verifyThingID(id)
	if err != nil {
		goto send
	}

	// Get the id generated as a token and send in the response
	token, err = i.thingProxy.SendCreateThing(id, name, authorization)

send:
	response := registerResponse{ID: id, Token: nil, Error: nil}
	if err != nil {
		response.Error = new(string)
		*response.Error = err.Error()
	}

	if len(token) > 0 {
		response.Token = new(string)
		*response.Token = token
	}

	bytes, err := json.Marshal(response)
	if err != nil {
		i.logger.Error(err)
		return err
	}

	fmt.Println(token)

	err = i.publisher.SendRegisterDevice(bytes)
	if err != nil {
		i.logger.Error(err)
		return err
	}

	return nil
}
