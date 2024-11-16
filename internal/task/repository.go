package task

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type repositoryDBImpl struct {
	db *gorm.DB
}

func NewRepositoryDB(db *gorm.DB) Repository {
	return &repositoryDBImpl{db: db}
}

func (r *repositoryDBImpl) Create(task *Task) error {
	model := &model{}
	model.FromEntity(task)

	if err := r.db.Create(model).Error; err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	task.ID = model.ID
	return nil
}

func (r *repositoryDBImpl) Get(id int) (*Task, error) {
	var model model
	if err := r.db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("task not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return model.ToEntity(), nil
}

func (r *repositoryDBImpl) Update(task *Task) error {
	model := &model{}
	model.FromEntity(task)

	if err := r.db.Save(model).Error; err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}
	return nil
}

func (r *repositoryDBImpl) Delete(id int) error {
	if err := r.db.Model(&model{}).Where("id = ?", id).
		Update("deleted_at", time.Now()).Error; err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

func (r *repositoryDBImpl) GetAll() ([]*Task, error) {
	var models []model
	if err := r.db.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get all tasks: %w", err)
	}

	tasks := make([]*Task, len(models))
	for i, model := range models {
		tasks[i] = model.ToEntity()
	}
	return tasks, nil
}

func (r *repositoryDBImpl) FindByCategory(category string) ([]*Task, error) {
	var models []model
	if err := r.db.Where("category = ?", category).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find tasks by category: %w", err)
	}

	tasks := make([]*Task, len(models))
	for i, model := range models {
		tasks[i] = model.ToEntity()
	}
	return tasks, nil
}

func (r *repositoryDBImpl) ExistsByMeta(metaType string, metaID int, category string) (bool, error) {
	var count int64
	err := r.db.Model(&model{}).
		Where("meta_type = ? AND meta_id = ? AND category = ? AND deleted_at IS NULL",
			metaType, metaID, category).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check task existence: %w", err)
	}
	return count > 0, nil
}
