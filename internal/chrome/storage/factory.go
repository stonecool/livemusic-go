package storage

import (
	"github.com/stonecool/livemusic-go/internal/chrome/instance"
	"github.com/stonecool/livemusic-go/internal/chrome/types"
)

func NewChrome(ip string, port int, url string, state types.ChromeState) *instance.Chrome {
	return instance.NewChrome(ip, port, url, state)
}
