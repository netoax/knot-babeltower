package interactors

import (
	"testing"

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

type FakeProxy struct {
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

func (fp *FakeProxy) SendCreateThing(id, name, authorization string) (idGenerated string, err error) {
	return idGenerated, nil
}

func TestRegisterThing(t *testing.T) {
	testCases := []struct {
		name          string
		testArguments bool
		thingID       string
		thingName     interface{}
		authorization interface{}
		fakeLogger    logging.Logger
		fakePublisher network.Publisher
		fakeProxy     network.ThingProxy
		errExpected   string
	}{
		{
			"shouldReturnNoError",
			false,
			"123",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			"",
		},
		{
			"shouldRaiseErrorIDLenght",
			false,
			"01234567890123456789",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			ErrorIDLenght{}.Error(),
		},
		{
			"shouldRaiseErrorNameNotFound",
			false,
			"123",
			"",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			ErrorNameNotFound{}.Error(),
		},
		{
			"shouldRaiseErrorIDInvalid",
			false,
			"not hex string",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			ErrorIDInvalid{}.Error(),
		},
		{
			"shouldRaisePublishError",
			false,
			"123",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisherWithSideEffect{},
			&FakeProxy{},
			ErrorFakePublisher{}.Error(),
		},
		{
			"shouldRaiseMissingArgument",
			true,
			"123",
			"",
			"",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			ErrorMissingArgument{}.Error(),
		},
		{
			"shouldInvalidTypeName",
			false,
			"123",
			123,
			"",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			ErrorInvalidTypeArgument{"Name is not string"}.Error(),
		},
		{
			"shouldInvalidTypeToken",
			false,
			"123",
			"test",
			123,
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			ErrorInvalidTypeArgument{"Authorization token is not string"}.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var err error
			createThingInteractor := NewRegisterThing(tc.fakeLogger, tc.fakePublisher, tc.fakeProxy)
			if tc.testArguments {
				err = createThingInteractor.Execute(tc.thingID)
			} else {
				err = createThingInteractor.Execute(tc.thingID, tc.thingName, tc.authorization)
			}

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
