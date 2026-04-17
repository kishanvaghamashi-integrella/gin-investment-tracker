package mocks

import "github.com/stretchr/testify/mock"

type MockPriceFetcher struct {
	mock.Mock
}

func (m *MockPriceFetcher) FetchPrice(instrumentType, externalPlatformID string) (float64, error) {
	args := m.Called(instrumentType, externalPlatformID)
	return args.Get(0).(float64), args.Error(1)
}
