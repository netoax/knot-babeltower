package mocks

import (
	"github.com/CESARBR/knot-babeltower/pkg/network"
	"github.com/CESARBR/knot-babeltower/pkg/thing/entities"
	"github.com/stretchr/testify/mock"
)

type FakePublisher struct {
	mock.Mock
	ReturnErr error
	SendError error
	Token     string
}

func (fp *FakePublisher) SendRegisterDevice(msg network.RegisterResponseMsg) error {
	ret := fp.Called(msg)
	return ret.Error(0)
}

func (fp *FakePublisher) SendUpdatedSchema(thingID string) error {
	ret := fp.Called(thingID)
	return ret.Error(0)
}

func (fp *FakePublisher) SendThings(things []*entities.Thing) error {
	args := fp.Called(things)
	return args.Error(0)
}
