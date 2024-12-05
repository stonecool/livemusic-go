package types

import (
	"testing"

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
				State:       0,
			},
			wantErr: false,
		},
		{
			name: "empty IP",
			model: &Model{
				IP:          "",
				Port:        9222,
				DebuggerURL: "ws://127.0.0.1:9222",
				State:       0,
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			model: &Model{
				IP:          "127.0.0.1",
				Port:        0,
				DebuggerURL: "ws://127.0.0.1:9222",
				State:       0,
			},
			wantErr: true,
		},
		{
			name: "empty debugger URL",
			model: &Model{
				IP:          "127.0.0.1",
				Port:        9222,
				DebuggerURL: "",
				State:       0,
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
	m := &Model{}
	assert.Equal(t, "chromes", m.TableName())
}

func TestModel_GetID(t *testing.T) {
	m := &Model{}
	m.ID = 123
	assert.Equal(t, 123, m.GetID())
}
