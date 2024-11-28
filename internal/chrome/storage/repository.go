package storage

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/chrome/instance"

	"github.com/stonecool/livemusic-go/internal/database"

	"gorm.io/gorm"
)

var (
	Repo repository
)

func init() {
	Repo = newRepositoryDB(database.DB)
}

type repository interface {
	Create(ip string, port int, debuggerURL string, state instance.ChromeState) (*instance.Chrome, error)
	Get(id int) (*instance.Chrome, error)
	Update(chrome *instance.Chrome) error
	Delete(id int) error
	GetAll() ([]*instance.Chrome, error)
	ExistsByIPAndPort(ip string, port int) (bool, error)
}

type repositoryDB struct {
	db database.Repository[*model]
}

func newRepositoryDB(db *gorm.DB) repository {
	return &repositoryDB{
		db: database.NewBaseRepository[*model](db),
	}
}

func (r *repositoryDB) Create(ip string, port int, debuggerURL string, state instance.ChromeState) (*instance.Chrome, error) {
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
		return nil, fmt.Errorf("failed to Create instance: %w", err)
	}

	return m.toEntity(), nil
}

func (r *repositoryDB) Get(id int) (*instance.Chrome, error) {
	m, err := r.db.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to Get instance: %w", err)
	}
	return m.toEntity(), nil
}

func (r *repositoryDB) Update(chrome *instance.Chrome) error {
	m := &model{}
	m.fromEntity(chrome)

	if err := r.db.Update(m).Error; err != nil {
		return fmt.Errorf("failed to Update instance: %w", err)
	}
	return nil
}

func (r *repositoryDB) Delete(id int) error {
	if err := r.db.Delete(id).Error; err != nil {
		return fmt.Errorf("failed to Delete instance: %w", err)
	}
	return nil
}

func (r *repositoryDB) GetAll() ([]*instance.Chrome, error) {
	models, err := r.db.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to Get task: %w", err)
	}

	var chromes []*instance.Chrome
	for _, m := range models {
		chromes = append(chromes, m.toEntity())
	}

	return chromes, nil
}

func (r *repositoryDB) ExistsByIPAndPort(ip string, port int) (bool, error) {
	return r.db.ExistsBy("ip = '?' AND port = '?'", ip, port)

}
