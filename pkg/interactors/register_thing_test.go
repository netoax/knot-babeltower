package interactors

import (
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/entities"
	"github.com/CESARBR/knot-babeltower/pkg/logging"
	"github.com/CESARBR/knot-babeltower/pkg/network"
)

type FakeRegisterThingLogger struct {
}

type FakeMsgPublisher struct {
}

type FakeMsgPublisherWithSideEffect struct {
}

type ErrorFakePublisher struct {
}

func (fl *FakeRegisterThingLogger) Info(...interface{}) {}

func (fl *FakeRegisterThingLogger) Infof(string, ...interface{}) {}

func (fl *FakeRegisterThingLogger) Debug(...interface{}) {}

func (fl *FakeRegisterThingLogger) Warn(...interface{}) {}

func (fl *FakeRegisterThingLogger) Error(...interface{}) {}

func (fl *FakeRegisterThingLogger) Errorf(string, ...interface{}) {}

func (em ErrorFakePublisher) Error() string {
	return "error mock"
}

func (fp *FakeMsgPublisher) SendRegisterDevice(msg []byte) error {
	return nil
}

func (fp *FakeMsgPublisherWithSideEffect) SendRegisterDevice(msg []byte) error {
	return ErrorFakePublisher{}
}

func TestRegisterThing(t *testing.T) {
	testCases := []struct {
		name          string
		thing         entities.Thing
		fakeLogger    logging.Logger
		fakePublisher network.Publisher
		errExpected   string
	}{
		{
			"shouldReturnNoError",
			entities.Thing{ID: "123", Name: "test"},
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			"",
		},
		{
			"shouldRaiseErrorIDLenght",
			entities.Thing{ID: "01234567890123456789", Name: "test"},
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			ErrorIDLenght{}.Error(),
		},
		{
			"shouldRaiseErrorNameNotFound",
			entities.Thing{ID: "123", Name: ""},
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			ErrorNameNotFound{}.Error(),
		},
		{
			"shouldRaiseErrorIDInvalid",
			entities.Thing{ID: "not hex string", Name: "test"},
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			ErrorIDInvalid{}.Error(),
		},
		{
			"shouldRaisePublishError",
			entities.Thing{ID: "123", Name: "test"},
			&FakeRegisterThingLogger{},
			&FakeMsgPublisherWithSideEffect{},
			ErrorFakePublisher{}.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			createThingInteractor := NewRegisterThing(tc.fakeLogger, tc.fakePublisher)
			err := createThingInteractor.Execute(tc.thing.ID, tc.thing.Name)
			if err != nil {
				if err.Error() != tc.errExpected {
					t.Errorf("Create Thing failed with unexpected error. Error: %s", err)
					return
				}

				t.Logf("Create Thing throws expected error. Error: %s", err)
			}

			t.Log("Create thing ok")
		})
	}
}
