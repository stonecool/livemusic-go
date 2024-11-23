package task

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/database"
	"gorm.io/gorm"
)

type repository interface {
	get(int) (*Task, error)
	create(category string, metaType string, metaID int, cronSpec string) (*Task, error)
	getAll() ([]*Task, error)
	existsByMeta(category string, metaType string, metaID int) (bool, error)
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

func (r *repositoryDB) create(category string, metaType string, metaID int, cronSpec string) (*Task, error) {
	exist, err := repo.existsByMeta(category, metaType, metaID)
	if !exist || err != nil {
		return nil, fmt.Errorf("exists")
	}

	model := &model{
		Category: category,
		MetaType: metaType,
		MetaID:   metaID,
		CronSpec: cronSpec,
	}

	if err := model.Validate(); err != nil {
		return nil, err
	}

	if err := r.db.Create(model); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return model.toEntity(), nil
}

func (r *repositoryDB) getAll() ([]*Task, error) {
	models, err := r.db.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	var tasks []*Task
	for _, m := range models {
		tasks = append(tasks, m.toEntity())
	}

	return tasks, nil
}

func (r *repositoryDB) existsByMeta(category string, metaType string, metaID int) (bool, error) {
	return r.db.ExistsBy("category = '?' AND metaType = '?' AND metaID = ?", category, metaType, metaID)
}
