package account

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccount_GetSetState(t *testing.T) {
	acc := &account{
		State: stateNew,
	}

	assert.Equal(t, stateNew, acc.getState())

	acc.SetState(stateInitialized)
	assert.Equal(t, stateInitialized, acc.getState())
}
