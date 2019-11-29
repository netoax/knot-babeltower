package interactors

import (
	"errors"
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type FakeUpdateSchemaLogger struct{}

type FakeThingProxy struct {
	mock.Mock
}

type FakePublisher struct {
	mock.Mock
}

type UpdateSchemaTestCase struct {
	name                   string
	thingID                string
	schemaList             []entities.Schema
	isSchemaValid          bool
	expectedGetThing       GetThingResponse
	expectedSchemaResponse error
	expectedUpdatedSchema  error
	fakeLogger             *FakeUpdateSchemaLogger
	fakeThingProxy         *FakeThingProxy
	fakePublisher          *FakePublisher
}

type GetThingResponse struct {
	err   error
	thing *entities.Thing
}

func (fl *FakeUpdateSchemaLogger) Info(...interface{}) {}

func (fl *FakeUpdateSchemaLogger) Infof(string, ...interface{}) {}

func (fl *FakeUpdateSchemaLogger) Debug(...interface{}) {}

func (fl *FakeUpdateSchemaLogger) Warn(...interface{}) {}

func (fl *FakeUpdateSchemaLogger) Error(...interface{}) {}

func (fl *FakeUpdateSchemaLogger) Errorf(string, ...interface{}) {}

func (ftp *FakeThingProxy) SendCreateThing(id, name, authorization string) (idGenerated string, err error) {
	return "", nil
}

func (ftp *FakeThingProxy) UpdateSchema(thingID string, schema []entities.Schema) error {
	args := ftp.Called(thingID, schema)
	return args.Error(0)
}

func (ftp *FakeThingProxy) Get(thingID string) (*entities.Thing, error) {
	args := ftp.Called(thingID)
	return args.Get(1).(*entities.Thing), args.Error(0)
}

func (fp *FakePublisher) SendRegisterDevice(msg []byte) error {
	return nil
}

func (fp *FakePublisher) SendUpdatedSchema(thingID string) error {
	args := fp.Called(thingID)
	return args.Error(0)
}

var cases = []UpdateSchemaTestCase{
	{
		"schema successfully updated on the thing's proxy",
		"29cf40c23012ce1c",
		[]entities.Schema{
			{
				SensorID:  0,
				ValueType: 3,
				Unit:      0,
				TypeID:    65521,
				Name:      "LED",
			},
		},
		true,
		GetThingResponse{nil, &entities.Thing{ID: "29cf40c23012ce1c"}},
		nil,
		nil,
		&FakeUpdateSchemaLogger{},
		&FakeThingProxy{},
		&FakePublisher{},
	},
	{
		"failed to update the schema on the thing's proxy",
		"29cf40c23012ce1c",
		[]entities.Schema{
			{
				SensorID:  0,
				ValueType: 3,
				Unit:      0,
				TypeID:    65521,
				Name:      "LED",
			},
		},
		true,
		GetThingResponse{nil, &entities.Thing{ID: "29cf40c23012ce1c"}},
		errors.New("failed to update schema"),
		nil,
		&FakeUpdateSchemaLogger{},
		&FakeThingProxy{},
		&FakePublisher{},
	},
	{
		"schema response successfully sent",
		"29cf40c23012ce1c",
		[]entities.Schema{
			{
				SensorID:  0,
				ValueType: 3,
				Unit:      0,
				TypeID:    65521,
				Name:      "LED",
			},
		},
		true,
		GetThingResponse{nil, &entities.Thing{ID: "29cf40c23012ce1c"}},
		nil,
		nil,
		&FakeUpdateSchemaLogger{},
		&FakeThingProxy{},
		&FakePublisher{},
	},
	{
		"failed to send updated schema response",
		"29cf40c23012ce1c",
		[]entities.Schema{
			{
				SensorID:  0,
				ValueType: 3,
				Unit:      0,
				TypeID:    65521,
				Name:      "LED",
			},
		},
		true,
		GetThingResponse{nil, &entities.Thing{ID: "29cf40c23012ce1c"}},
		nil,
		errors.New("failed to send updated schema response"),
		&FakeUpdateSchemaLogger{},
		&FakeThingProxy{},
		&FakePublisher{},
	},
	{
		"thing doesn't exist",
		"29cf40c23012ce1c",
		[]entities.Schema{
			{
				SensorID:  0,
				ValueType: 3,
				Unit:      0,
				TypeID:    65521,
				Name:      "LED",
			},
		},
		true,
		GetThingResponse{errors.New("thing doesn't exist"), nil},
		nil,
		nil,
		&FakeUpdateSchemaLogger{},
		&FakeThingProxy{},
		&FakePublisher{},
	},
	{
		"invalid schema",
		"29cf40c23012ce1c",
		[]entities.Schema{
			{
				SensorID:  0,
				ValueType: 3,
				Unit:      0,
				TypeID:    65521,
				Name:      "LED",
			},
		},
		false,
		GetThingResponse{nil, &entities.Thing{ID: "29cf40c23012ce1c"}},
		nil,
		nil,
		&FakeUpdateSchemaLogger{},
		&FakeThingProxy{},
		&FakePublisher{},
	},
}

func TestUpdateSchema(t *testing.T) {
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.fakeThingProxy.
				On("Get", tc.thingID).
				Return(tc.expectedGetThing.err, tc.expectedGetThing.thing)

			tc.fakeThingProxy.
				On("UpdateSchema", tc.thingID, tc.schemaList).
				Return(tc.expectedSchemaResponse).
				Maybe()

			tc.fakePublisher.
				On("SendUpdatedSchema", tc.thingID).
				Return(tc.expectedUpdatedSchema).
				Maybe()

			updateSchemaInteractor := NewUpdateSchema(tc.fakeLogger, tc.fakePublisher, tc.fakeThingProxy)
			err := updateSchemaInteractor.Execute(tc.thingID, tc.schemaList)
			if !tc.isSchemaValid {
				assert.EqualError(t, err, "Thing's schema is invalid")
			}

			tc.fakeThingProxy.AssertExpectations(t)
			tc.fakePublisher.AssertExpectations(t)
		})
	}
}

/*
	1.	check if device exists
	2.	validate schema
	3. 	update schema
	4.	send updated schema response
	... add logger
	5.  send schema message to connector
*/
