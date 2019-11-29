package interactors

import (
	"fmt"

	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/go-playground/validator"
)

// // UpdateSchema represents the update schema use case
// type UpdateSchema struct {
// 	logger     logging.Logger
// 	publisher  amqp.Publisher
// 	thingProxy http.ThingProxy
// }

// ErrInvalidSchema represents the error when the schema has a invalid format
type ErrInvalidSchema struct{}

func (eis *ErrInvalidSchema) Error() string {
	return "Thing's schema is invalid"
}

// ErrThingExists represents the error when the schema has a invalid format
type ErrThingExists struct {
	thingID string
}

func (ete *ErrThingExists) Error() string {
	return fmt.Sprintf("Thing %s already exists", ete.thingID)
}

// // NewUpdateSchema creates a new UpdateSchema interactor instance
// func NewUpdateSchema(logger logging.Logger, publisher amqp.Publisher, thingProxy http.ThingProxy) UpdateSchema {
// 	return UpdateSchema{logger, publisher, thingProxy}
// }

// UpdateSchema receive the new sensor schema and update it on the thing's service
func (i *ThingInteractor) UpdateSchema(authorization, thingID string, schemaList []entities.Schema) error {
	if !i.thingExists(thingID) {
		return &ErrThingExists{thingID}
	}

	if !i.isValidSchema(schemaList) {
		return &ErrInvalidSchema{}
	}

	err := i.thingProxy.UpdateSchema(thingID, schemaList)
	if err != nil {
		return err
	}

	err = i.publisher.SendUpdatedSchema(thingID)
	if err != nil {
		return err
	}

	return nil
}

func (i *ThingInteractor) isValidSchema(schemaList []entities.Schema) bool {
	validate := validator.New()
	for _, schema := range schemaList {
		err := validate.Struct(schema)
		if err != nil {
			return false
		}
	}

	return true
}

func (i *ThingInteractor) thingExists(thingID string) bool {
	thing, err := i.thingProxy.Get(thingID)
	if err != nil {
		return false
	}

	if thing != nil {
		return true
	}

	return false
}
