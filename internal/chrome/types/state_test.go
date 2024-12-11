package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstanceState_String(t *testing.T) {
	tests := []struct {
		name  string
		state InstanceState
		want  string
	}{
		{
			name:  "Invalid state",
			state: InstanceStateInvalid,
			want:  "Invalid",
		},
		{
			name:  "Available state",
			state: InstanceStateAvailable,
			want:  "Available",
		},
		{
			name:  "Unstable state",
			state: InstanceStateUnstable,
			want:  "Unstable",
		},
		{
			name:  "Unavailable state",
			state: InstanceStateUnavailable,
			want:  "Unavailable",
		},
		{
			name:  "Unknown state",
			state: InstanceState(99),
			want:  "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.state.String())
		})
	}
}

func TestInstanceState_IsValidTransition(t *testing.T) {
	tests := []struct {
		name      string
		state     InstanceState
		event     EventType
		wantValid bool
	}{
		{
			name:      "Invalid state cannot transition",
			state:     InstanceStateInvalid,
			event:     EventHealthCheckSuccess,
			wantValid: false,
		},
		{
			name:      "Available can become unstable",
			state:     InstanceStateAvailable,
			event:     EventHealthCheckFail,
			wantValid: true,
		},
		{
			name:      "Available can stay available",
			state:     InstanceStateAvailable,
			event:     EventHealthCheckSuccess,
			wantValid: true,
		},
		{
			name:      "Unstable can recover",
			state:     InstanceStateUnstable,
			event:     EventHealthCheckSuccess,
			wantValid: true,
		},
		{
			name:      "Unstable can become unavailable",
			state:     InstanceStateUnstable,
			event:     EventHealthCheckFail,
			wantValid: true,
		},
		{
			name:      "Unavailable can recover",
			state:     InstanceStateUnavailable,
			event:     EventHealthCheckSuccess,
			wantValid: true,
		},
		{
			name:      "Unavailable cannot become unstable",
			state:     InstanceStateUnavailable,
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

func TestValidTransitions(t *testing.T) {
	// Test that all states have defined transitions
	states := []InstanceState{
		InstanceStateInvalid,
		InstanceStateAvailable,
		InstanceStateUnstable,
		InstanceStateUnavailable,
	}

	for _, state := range states {
		t.Run(state.String(), func(t *testing.T) {
			transitions, exists := validTransitions[state]
			assert.True(t, exists, "State should have defined transitions")
			assert.NotNil(t, transitions, "Transitions should not be nil")
		})
	}
}
