package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChromeState_String(t *testing.T) {
	tests := []struct {
		name  string
		state ChromeState
		want  string
	}{
		{
			name:  "connected state",
			state: ChromeStateConnected,
			want:  "Connected",
		},
		{
			name:  "disconnected state",
			state: ChromeStateDisconnected,
			want:  "Disconnected",
		},
		{
			name:  "offline state",
			state: ChromeStateOffline,
			want:  "Offline",
		},
		{
			name:  "unknown state",
			state: ChromeState(99),
			want:  "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.state.String())
		})
	}
}

func TestChromeState_IsValidTransition(t *testing.T) {
	tests := []struct {
		name      string
		state     ChromeState
		event     EventType
		wantValid bool
	}{
		{
			name:      "connected to disconnected",
			state:     ChromeStateConnected,
			event:     EventHealthCheckFail,
			wantValid: true,
		},
		{
			name:      "invalid transition",
			state:     ChromeStateConnected,
			event:     EventShutdown,
			wantValid: false,
		},
		{
			name:      "unknown state",
			state:     ChromeState(99),
			event:     EventHealthCheckFail,
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantValid, tt.state.IsValidTransition(tt.event))
		})
	}
}
