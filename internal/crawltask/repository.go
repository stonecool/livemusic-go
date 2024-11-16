package crawltask

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(task *CrawlTask) error
	Get(id int) (*CrawlTask, error)
	Update(task *CrawlTask) error
	Delete(id int) error
	GetAll() ([]*CrawlTask, error)
	FindByCategory(category string) ([]*CrawlTask, error)
	ExistsByMeta(metaType string, metaID int, category string) (bool, error)
}

type repositoryDBImpl struct {
	db *gorm.DB
}

func NewRepositoryDB(db *gorm.DB) Repository {
	return &repositoryDBImpl{db: db}
}

func (r *repositoryDBImpl) Create(task *CrawlTask) error {
	model := &Model{}
	model.FromEntity(task)

	if err := r.db.Create(model).Error; err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	task.ID = model.ID
	return nil
}

func (r *repositoryDBImpl) Get(id int) (*CrawlTask, error) {
	var model Model
	if err := r.db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("task not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return model.ToEntity(), nil
}

func (r *repositoryDBImpl) Update(task *CrawlTask) error {
	model := &Model{}
	model.FromEntity(task)

	if err := r.db.Save(model).Error; err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	return nil
}

func (r *repositoryDBImpl) Delete(id int) error {
	if err := r.db.Model(&Model{}).Where("id = ?", id).
		Update("deleted_at", time.Now()).Error; err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

func (r *repositoryDBImpl) GetAll() ([]*CrawlTask, error) {
	var models []Model
	if err := r.db.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get all tasks: %w", err)
	}

	tasks := make([]*CrawlTask, len(models))
	for i, model := range models {
		tasks[i] = model.ToEntity()
	}
	return tasks, nil
}

func (r *repositoryDBImpl) FindByCategory(category string) ([]*CrawlTask, error) {
	var models []Model
	if err := r.db.Where("category = ?", category).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find tasks by category: %w", err)
	}

	tasks := make([]*CrawlTask, len(models))
	for i, model := range models {
		tasks[i] = model.ToEntity()
	}
	return tasks, nil
}

func (r *repositoryDBImpl) ExistsByMeta(metaType string, metaID int, category string) (bool, error) {
	var count int64
	err := r.db.Model(&Model{}).
		Where("meta_type = ? AND meta_id = ? AND category = ? AND deleted_at IS NULL",
			metaType, metaID, category).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check task existence: %w", err)
	}
	return count > 0, nil
}
