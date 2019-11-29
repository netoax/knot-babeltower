package interactors

import (
	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/thing/delivery/amqp"
	"github.com/CESARBR/knot-babeltower/pkg/thing/delivery/http"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
)

// Interactor is an interface that defines the thing's use cases operations
type Interactor interface {
	Register(authorization, id, name string) error
	UpdateSchema(authorization, id string, schemaList []entities.Schema) error
}

// ThingInteractor represents the thing's interactor
type ThingInteractor struct {
	logger     logging.Logger
	publisher  amqp.Publisher
	thingProxy http.ThingProxy
}

// NewThingInteractor contructs the use case
func NewThingInteractor(logger logging.Logger, publisher amqp.Publisher, thingProxy http.ThingProxy) *ThingInteractor {
	return &ThingInteractor{logger, publisher, thingProxy}
}
