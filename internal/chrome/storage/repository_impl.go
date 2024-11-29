package storage

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/chrome/types"

	"github.com/stonecool/livemusic-go/internal/chrome/instance"
	"github.com/stonecool/livemusic-go/internal/database"
	"gorm.io/gorm"
)

var (
	Repo Repository
)

func init() {
	Repo = newRepositoryDB(database.DB)
}

type repositoryDB struct {
	db database.Repository[*model]
}

func newRepositoryDB(db *gorm.DB) Repository {
	return &repositoryDB{
		db: database.NewBaseRepository[*model](db),
	}
}

func (r *repositoryDB) Create(ip string, port int, debuggerURL string, state types.ChromeState) (types.IChrome, error) {
	m := &model{
		IP:          ip,
		Port:        port,
		DebuggerURL: debuggerURL,
		State:       int(state),
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}

	if err := r.db.Create(m); err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	return m.toEntity(), nil
}

func (r *repositoryDB) Get(id int) (types.IChrome, error) {
	m, err := r.db.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}
	return m.toEntity(), nil
}

func (r *repositoryDB) Update(chrome types.IChrome) error {
	// 由于接口转换的问题，这里需要做一个类型断言
	chromeInstance, ok := chrome.(*instance.Chrome)
	if !ok {
		return fmt.Errorf("invalid chrome instance type")
	}

	m := &model{}
	m.fromEntity(chromeInstance)

	if err := r.db.Update(m).Error; err != nil {
		return fmt.Errorf("failed to update instance: %w", err)
	}
	return nil
}

func (r *repositoryDB) Delete(id int) error {
	if err := r.db.Delete(id).Error; err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}
	return nil
}

func (r *repositoryDB) GetAll() ([]types.IChrome, error) {
	models, err := r.db.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all instances: %w", err)
	}

	chromes := make([]types.IChrome, len(models))
	for i, m := range models {
		chromes[i] = m.toEntity()
	}

	return chromes, nil
}

func (r *repositoryDB) ExistsByIPAndPort(ip string, port int) (bool, error) {
	return r.db.ExistsBy("ip = ? AND port = ?", ip, port)
}
