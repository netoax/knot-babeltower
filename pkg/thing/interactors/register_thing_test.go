package interactors

import (
	"testing"

	"github.com/CESARBR/knot-babeltower/pkg/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type FakeRegisterThingLogger struct {
}

type FakeMsgPublisher struct {
}

type ErrorFakePublisher struct {
}

type FakeProxy struct {
	mock.Mock
}

type ErrorFakeProxy struct {
}

func (fl *FakeRegisterThingLogger) Info(...interface{}) {}

func (fl *FakeRegisterThingLogger) Infof(string, ...interface{}) {}

func (fl *FakeRegisterThingLogger) Debug(...interface{}) {}

func (fl *FakeRegisterThingLogger) Warn(...interface{}) {}

func (fl *FakeRegisterThingLogger) Error(...interface{}) {}

func (fl *FakeRegisterThingLogger) Errorf(string, ...interface{}) {}

func (em ErrorFakePublisher) Error() string {
	return "error publish mock"
}

func (fp *FakeMsgPublisher) SendRegisterDevice(msg []byte) error {
	return nil
}

func (fp *FakeMsgPublisher) SendUpdatedSchema(thingID string) error {
	return nil
}

func (fp *FakeProxy) SendCreateThing(id, name, authorization string) (idGenerated string, err error) {
	ret := fp.Called(id, name, authorization)

	rf, ok := ret.Get(0).(func(string, string, string) (string, error))
	if ok {
		idGenerated, err = rf(id, name, authorization)
	} else {
		idGenerated, err = ret.String(0), ret.Error(1)
	}

	return idGenerated, err
}

func (fp *FakeProxy) UpdateSchema(ID string, schemaList []entities.Schema) error {
	return nil
}

func (fp *FakeProxy) Get(ID string) (*entities.Thing, error) {
	return nil, nil
}

func (em ErrorFakeProxy) Error() string {
	return "error proxy mock"
}

func TestRegisterThing(t *testing.T) {
	testCases := map[string]struct {
		testArguments    bool
		thingID          string
		thingName        interface{}
		authorization    interface{}
		fakeLogger       *FakeRegisterThingLogger
		fakePublisher    *FakeMsgPublisher
		fakeProxy        *FakeProxy
		proxyReturnID    string
		proxyReturnError error
		errExpected      error
	}{
		"shouldReturnNoError": {
			false,
			"123",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			"secret",
			nil,
			nil,
		},
		"TestIDLenght": {
			false,
			"01234567890123456789",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			"secret",
			nil,
			nil,
		},
		"TestNameEmpty": {
			false,
			"123",
			"",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			"secret",
			nil,
			ErrorNameNotFound{},
		},
		"TestIDInvalid": {
			false,
			"not hex string",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			"secret",
			nil,
			ErrorIDInvalid{},
		},
		"shouldRaisePublishError": {
			false,
			"123",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			"secret",
			nil,
			ErrorFakePublisher{},
		},
		"TestProxyError": {
			false,
			"123",
			"test",
			"authorization token",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			"secret",
			ErrorFakeProxy{},
			nil,
		},
		"shouldRaiseMissingAuthorizationToken": {
			false,
			"123",
			"test",
			"",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			"secret",
			nil,
			ErrorUnauthorized{},
		},
		"shouldRaiseMissingArgument": {
			true,
			"123",
			"",
			"",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			"secret",
			nil,
			ErrorMissingArgument{},
		},
		"shouldInvalidTypeName": {
			false,
			"123",
			123,
			"",
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			"secret",
			nil,
			ErrorInvalidTypeArgument{"Name is not string"},
		},
		"shouldInvalidTypeToken": {
			false,
			"123",
			"test",
			123,
			&FakeRegisterThingLogger{},
			&FakeMsgPublisher{},
			&FakeProxy{},
			"secret",
			nil,
			ErrorInvalidTypeArgument{"Authorization token is not string"},
		},
	}

	t.Logf("Number of test cases: %d", len(testCases))
	for tcName, tc := range testCases {
		t.Logf("Test case %s", tcName)
		t.Run(tcName, func(t *testing.T) {
			var err error
			tc.fakeProxy.On("SendCreateThing", tc.thingID, tc.thingName, tc.authorization).
				Return(tc.proxyReturnID, tc.proxyReturnError).Maybe()
			createThingInteractor := NewRegisterThing(tc.fakeLogger, tc.fakePublisher, tc.fakeProxy)
			if tc.testArguments {
				err = createThingInteractor.Execute(tc.thingID)
			} else {
				err = createThingInteractor.Execute(tc.thingID, tc.thingName, tc.authorization)
			}

			if err != nil && !assert.IsType(t, err, tc.errExpected) {
				t.Errorf("Create Thing failed with unexpected error. Error: %s", err)
				return
			}

			t.Log("Create thing ok")
			tc.fakeProxy.AssertExpectations(t)
		})
	}
}
