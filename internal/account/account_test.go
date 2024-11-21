package account

import (
	"testing"

	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/message"
	"github.com/stretchr/testify/assert"
)

func TestAccount_Init(t *testing.T) {
	acc := &Account{
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
	acc := &Account{
		State: internal.AS_EXPIRED,
	}

	assert.Equal(t, internal.AS_EXPIRED, acc.GetState())

	acc.SetState(internal.AS_RUNNING)
	assert.Equal(t, internal.AS_RUNNING, acc.GetState())
}
