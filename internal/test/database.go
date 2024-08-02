package test

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockCommitStore struct {
	mock.Mock
}

func (m *MockCommitStore) GetLastCommitDate() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}
