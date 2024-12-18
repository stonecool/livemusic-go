package account

import (
	"github.com/stonecool/livemusic-go/internal/account/state"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccount_GetSetState(t *testing.T) {
	acc := &account{
		State: state.stateNew,
	}

	assert.Equal(t, state.stateNew, acc.getState())

	acc.SetState(stateInitialized)
	assert.Equal(t, stateInitialized, acc.getState())
}
