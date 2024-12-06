package storage

import "github.com/stonecool/livemusic-go/internal/chrome/types"

type Repository interface {
	Create(ip string, port int, debuggerURL string, state types.ChromeState) (*types.Model, error)
	Get(int) (*types.Model, error)
	Update(model *types.Model) error
	Delete(int) error
	GetAll() ([]*types.Model, error)
	ExistsByIPAndPort(ip string, port int) (bool, error)
}
