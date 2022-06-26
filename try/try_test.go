package try

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type runMocker struct {
	mock.Mock
}

func (r *runMocker) run() error {
	return r.Called().Error(0)
}

func before(m *runMocker, nums, failAt int) []Try {
	res := make([]Try, nums)
	for i := range res {
		var err error
		if failAt == i {
			err = fmt.Errorf("intended at %d", failAt)
		}
		m.On("run").Return(err).Once()
		res[i] = func() error {
			return m.run()
		}
	}
	return res
}

func TestRunOrError(t *testing.T) {
	testCases := []struct {
		nums       int
		shouldFail bool
		failAt     int
	}{
		{1, true, 0},
		{2, true, 1},
		{2, false, -1},
		{10, false, -1},
		{10, true, 5},
	}
	for _, c := range testCases {
		m := &runMocker{}
		tries := before(m, c.nums, c.failAt)
		err := RunOrError(tries...)
		if c.shouldFail {
			assert.NotNil(t, err)
			m.AssertNumberOfCalls(t, "run", c.failAt+1)
			continue
		}
		assert.Nil(t, err)
		m.AssertNumberOfCalls(t, "run", c.nums)
	}
}
