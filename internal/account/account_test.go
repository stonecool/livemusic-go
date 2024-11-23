package account

import (
	"testing"

	"github.com/stonecool/livemusic-go/internal/message"
	"github.com/stretchr/testify/assert"
)

func TestAccount_Init(t *testing.T) {
	account := &account{
		ID:      1,
		msgChan: make(chan *message.AsyncMessage),
		done:    make(chan struct{}),
	}

	account.Init()
	defer account.Close()

	assert.NotNil(t, account.msgChan)
	assert.NotNil(t, account.done)
}

func TestAccount_GetSetState(t *testing.T) {
	account := &account{
		State: stateNew,
	}

	assert.Equal(t, stateNew, account.getState())

	account.SetState(stateInitialized)
	assert.Equal(t, stateInitialized, account.getState())
}
