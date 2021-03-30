package press

import (
	"errors"
	"math/rand"
	"testing"
	"time"
)

type TestAgent struct {
	r *rand.Rand
}

func NewTestAgent() Agent {
	return &TestAgent{
		r: rand.New(rand.NewSource(100)),
	}
}

func (s *TestAgent) Execute() error {
	r := s.r.Int() % 10
	time.Sleep(time.Millisecond * time.Duration(r))
	if r >= 5 {
		return errors.New("test error")
	} else {
		return nil
	}
}

func TestNewPressController(t *testing.T) {
	controller := NewController(NewTestAgent, &Option{Number: 100})
	controller.Start()
}
