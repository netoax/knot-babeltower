package mocks

type FakeLogger struct{}

func (fl *FakeLogger) Info(...interface{}) {}

func (fl *FakeLogger) Infof(string, ...interface{}) {}

func (fl *FakeLogger) Debug(...interface{}) {}

func (fl *FakeLogger) Warn(...interface{}) {}

func (fl *FakeLogger) Error(...interface{}) {}

func (fl *FakeLogger) Errorf(string, ...interface{}) {}
