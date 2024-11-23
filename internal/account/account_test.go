package account

import (
	"testing"

	"github.com/stonecool/livemusic-go/internal/message"
	"github.com/stretchr/testify/assert"
)

func TestAccount_Init(t *testing.T) {
	acc := &account{
		ID:      1,
		msgChan: make(chan *message.AsyncMessage),
		done:    make(chan struct{}),
	}

	acc.Init()
	defer acc.Close()

	assert.NotNil(t, acc.msgChan)
	assert.NotNil(t, acc.done)
}

func TestAccount_GetSetState(t *testing.T) {
	acc := &account{
		State: stateNew,
	}

	assert.Equal(t, stateNew, acc.getState())

	acc.SetState(stateInitialized)
	assert.Equal(t, stateInitialized, acc.getState())
}
