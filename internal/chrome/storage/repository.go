package storage

import "github.com/stonecool/livemusic-go/internal/chrome/types"

type Repository interface {
	Get(id int) (*types.Model, error)
	Create(ip string, port int, debuggerURL string, state types.InstanceState) (*types.Model, error)
	Update(model *types.Model) error
	Delete(id int) error
	GetAll() ([]*types.Model, error)
	ExistsByIPAndPort(ip string, port int) (bool, error)
}
