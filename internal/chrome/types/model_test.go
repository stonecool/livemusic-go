package types

import (
	"testing"

	"github.com/stonecool/livemusic-go/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestModel_Validate(t *testing.T) {
	tests := []struct {
		name    string
		model   *Model
		wantErr bool
	}{
		{
			name: "valid model",
			model: &Model{
				IP:          "127.0.0.1",
				Port:        9222,
				DebuggerURL: "ws://127.0.0.1:9222",
				State:       int(InstanceStateAvailable),
			},
			wantErr: false,
		},
		{
			name: "empty IP",
			model: &Model{
				IP:          "",
				Port:        9222,
				DebuggerURL: "ws://127.0.0.1:9222",
				State:       int(InstanceStateAvailable),
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			model: &Model{
				IP:          "127.0.0.1",
				Port:        0,
				DebuggerURL: "ws://127.0.0.1:9222",
				State:       int(InstanceStateAvailable),
			},
			wantErr: true,
		},
		{
			name: "empty debugger URL",
			model: &Model{
				IP:          "127.0.0.1",
				Port:        9222,
				DebuggerURL: "",
				State:       int(InstanceStateAvailable),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestModel_TableName(t *testing.T) {
	model := Model{}
	assert.Equal(t, "chromes", model.TableName())
}

func TestModel_GetID(t *testing.T) {
	m := &Model{}
	m.ID = 123
	assert.Equal(t, 123, m.GetID())
}

func TestModel_StateHandling(t *testing.T) {
	tests := []struct {
		name  string
		state int
	}{
		{
			name:  "Invalid state",
			state: int(InstanceStateInvalid),
		},
		{
			name:  "Available state",
			state: int(InstanceStateAvailable),
		},
		{
			name:  "Unstable state",
			state: int(InstanceStateUnstable),
		},
		{
			name:  "Unavailable state",
			state: int(InstanceStateUnavailable),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := &Model{State: tt.state}
			assert.Equal(t, tt.state, model.State)
		})
	}
}

func TestModel_Fields(t *testing.T) {
	model := &Model{
		BaseModel: database.BaseModel{
			ID: 123,
		},
		IP:          "127.0.0.1",
		Port:        9222,
		DebuggerURL: "ws://127.0.0.1:9222",
		State:       int(InstanceStateAvailable),
	}

	assert.Equal(t, 123, model.ID)
	assert.Equal(t, "127.0.0.1", model.IP)
	assert.Equal(t, 9222, model.Port)
	assert.Equal(t, "ws://127.0.0.1:9222", model.DebuggerURL)
	assert.Equal(t, int(InstanceStateAvailable), model.State)
}
