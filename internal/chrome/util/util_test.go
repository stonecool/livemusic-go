package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidIPv4(t *testing.T) {
	tests := []struct {
		ip       string
		expected bool
	}{
		{"192.168.1.1", true},
		{"255.255.255.255", true},
		{"0.0.0.0", true},
		{"256.256.256.256", false},
		{"abc.def.ghi.jkl", false},
		{"", false},
	}

	for _, test := range tests {
		result := IsValidIPv4(test.ip)
		if result != test.expected {
			t.Errorf("IsValidIPv4(%s) = %v; want %v", test.ip, result, test.expected)
		}
	}
}

func TestIsValidPort(t *testing.T) {
	tests := []struct {
		port     int
		expected bool
	}{
		{80, true},
		{0, true},
		{65535, true},
		{-1, false},
		{65536, false},
	}

	for _, test := range tests {
		result := IsValidPort(test.port)
		if result != test.expected {
			t.Errorf("IsValidPort(%d) = %v; want %v", test.port, result, test.expected)
		}
	}
}

func TestConvertIPv4ToInt(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		want    uint32
		wantErr bool
	}{
		{
			name:    "localhost",
			ip:      "127.0.0.1",
			want:    2130706433,
			wantErr: false,
		},
		{
			name:    "invalid format",
			ip:      "127.0.0",
			wantErr: true,
		},
		{
			name:    "invalid number",
			ip:      "127.0.0.256",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertIPv4ToInt(tt.ip)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
