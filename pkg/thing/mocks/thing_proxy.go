package mocks

import (
	"github.com/CESARBR/knot-babeltower/pkg/thing/delivery/http"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/stretchr/testify/mock"
)

type FakeThingProxy struct {
	mock.Mock
	ReturnErr error
}

func (ftp *FakeThingProxy) Create(id, name, authorization string) (idGenerated string, err error) {
	ret := ftp.Called(id, name, authorization)
	return ret.String(0), ret.Error(1)
}

func (ftp *FakeThingProxy) UpdateSchema(authorization, thingID string, schema []entities.Schema) error {
	ret := ftp.Called(thingID, schema)
	return ret.Error(0)
}

func (ftp *FakeThingProxy) Get(authorization, thingID string) (*http.ThingProxyRepr, error) {

	return nil, nil
}

func (ftp *FakeThingProxy) List(authorization string) ([]*entities.Thing, error) {
	// args := ftp.Called(authorization)
	// return args.Get(0), args.Error(1)
	return nil, nil
}
