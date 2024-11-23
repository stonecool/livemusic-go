package task

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/database"
	"gorm.io/gorm"
)

type repository interface {
	get(int) (*Task, error)
	create(category string, metaType string, metaId int, cronSpec string) (*Task, error)
	getAll() ([]*Task, error)
}

type repositoryDB struct {
	db database.Repository[*model]
}

func newRepositoryDB(db *gorm.DB) repository {
	return &repositoryDB{
		db: database.NewBaseRepository[*model](db),
	}
}

func (r *repositoryDB) get(id int) (*Task, error) {
	m, err := r.db.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	return m.toEntity(), nil
}

func (r *repositoryDB) create(category string, metaType string, metaId int, cronSpec string) (*Task, error) {
	m := &model{
		Category: category,
		MetaType: metaType,
		MetaId:   metaId,
		CronSpec: cronSpec,
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}

	if err := r.db.Create(m); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return m.toEntity(), nil
}

func (r *repositoryDB) getAll() ([]*Task, error) {
	return nil, nil
}
