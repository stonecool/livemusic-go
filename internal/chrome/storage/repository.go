package storage

import (
	"github.com/stonecool/livemusic-go/internal/chrome/types"
)

type Repository interface {
	Create(ip string, port int, debuggerURL string, state types.ChromeState) (types.IChrome, error)
	Get(id int) (types.IChrome, error)
	Update(chrome types.IChrome) error
	Delete(id int) error
	GetAll() ([]types.IChrome, error)
	ExistsByIPAndPort(ip string, port int) (bool, error)
}
