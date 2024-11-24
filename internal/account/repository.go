package account

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/database"
	"gorm.io/gorm"
)

type repository interface {
	get(int) (*account, error)
	create(string, state) (*account, error)
}

type repositoryDB struct {
	db database.Repository[*model]
}

func newRepositoryDB(db *gorm.DB) repository {
	return &repositoryDB{
		db: database.NewBaseRepository[*model](db),
	}
}

func (r *repositoryDB) create(category string, state state) (*account, error) {
	m := &model{
		Category: category,
		State:    int(state),
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}

	if err := r.db.Create(m); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return m.toEntity(), nil
}

func (r *repositoryDB) get(id int) (*account, error) {
	m, err := r.db.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return m.toEntity(), nil
}
