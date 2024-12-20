package account

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/database"
	"github.com/stonecool/livemusic-go/internal/message"
	"gorm.io/gorm"
)

var (
	repo repository
)

func init() {
	repo = newRepositoryDB(database.DB)
}

type repository interface {
	create(string, message.AccountState) (*account, error)
	get(int) (*account, error)
}

type repositoryDB struct {
	db database.Repository[*model]
}

func newRepositoryDB(db *gorm.DB) repository {
	return &repositoryDB{
		db: database.NewBaseRepository[*model](db),
	}
}

func (r *repositoryDB) create(category string, state message.AccountState) (*account, error) {
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
