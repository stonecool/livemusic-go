package account

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/database"
	"gorm.io/gorm"
)

type repository interface {
	get(int) (*Account, error)
	create(string, state) (*Account, error)
}

type repositoryDB struct {
	db database.Repository[*accountModel]
}

func newRepositoryDB(db *gorm.DB) repository {
	return &repositoryDB{
		db: database.NewBaseRepository[*accountModel](db),
	}
}

func (r *repositoryDB) get(id int) (*Account, error) {
	m, err := r.db.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return m.toEntity(), nil
}

func (r *repositoryDB) create(category string, s state) (*Account, error) {
	m := &accountModel{
		Category: category,
		State:    int(s),
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}

	if err := r.db.Create(m); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return m.toEntity(), nil
}
