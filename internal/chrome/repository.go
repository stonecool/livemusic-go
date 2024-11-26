package chrome

import (
	"fmt"

	"github.com/stonecool/livemusic-go/internal/database"

	"gorm.io/gorm"
)

type repository interface {
	create(ip string, port int, debuggerURL string, state State) (*Chrome, error)
	get(id int) (*Chrome, error)
	update(chrome *Chrome) error
	delete(id int) error
	getAll() ([]*Chrome, error)
	existsByIPAndPort(ip string, port int) (bool, error)
}

type repositoryDB struct {
	db database.Repository[*model]
}

func newRepositoryDB(db *gorm.DB) repository {
	return &repositoryDB{
		db: database.NewBaseRepository[*model](db),
	}
}

func (r *repositoryDB) create(ip string, port int, debuggerURL string, state State) (*Chrome, error) {
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

func (r *repositoryDB) get(id int) (*Chrome, error) {
	m, err := r.db.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}
	return m.toEntity(), nil
}

func (r *repositoryDB) update(chrome *Chrome) error {
	m := &model{}
	m.fromEntity(chrome)

	if err := r.db.Update(m).Error; err != nil {
		return fmt.Errorf("failed to update instance: %w", err)
	}
	return nil
}

func (r *repositoryDB) delete(id int) error {
	if err := r.db.Delete(id).Error; err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}
	return nil
}

func (r *repositoryDB) getAll() ([]*Chrome, error) {
	models, err := r.db.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	var chromes []*Chrome
	for _, m := range models {
		chromes = append(chromes, m.toEntity())
	}

	return chromes, nil
}

func (r *repositoryDB) existsByIPAndPort(ip string, port int) (bool, error) {
	return r.db.ExistsBy("ip = '?' AND port = '?'", ip, port)

}
