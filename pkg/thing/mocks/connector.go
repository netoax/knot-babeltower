package mocks

import "github.com/stretchr/testify/mock"

type FakeConnector struct {
	mock.Mock
	SendError error
	RecvError error
}

func (fc *FakeConnector) SendRegisterDevice(id, name string) (err error) {
	ret := fc.Called(id, name)
	return ret.Error(0)
}

func (fc *FakeConnector) RecvRegisterDevice() (bytes []byte, err error) {
	ret := fc.Called()
	return bytes, ret.Error(1)
}
