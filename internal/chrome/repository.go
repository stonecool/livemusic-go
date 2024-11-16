package chrome

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	Create(instance *Instance) error
	Get(id int) (*Instance, error)
	Update(instance *Instance) error
	Delete(id int) error
	GetAll() ([]*Instance, error)
	FindByCategory(category string) ([]*Instance, error)
}

type repositoryDBImpl struct {
	db *gorm.DB
}

func NewRepositoryDB(db *gorm.DB) Repository {
	return &repositoryDBImpl{db: db}
}

func (r *repositoryDBImpl) Create(instance *Instance) error {
	model := &Model{}
	model.FromEntity(instance)

	if err := r.db.Create(model).Error; err != nil {
		return fmt.Errorf("failed to create instance: %w", err)
	}

	instance.ID = model.ID
	return nil
}

func (r *repositoryDBImpl) Get(id int) (*Instance, error) {
	var model Model
	if err := r.db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("instance not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}

	instance := model.ToEntity()
	return instance, nil
}

func (r *repositoryDBImpl) Update(instance *Instance) error {
	model := &Model{}
	model.FromEntity(instance)

	if err := r.db.Save(model).Error; err != nil {
		return fmt.Errorf("failed to update instance: %w", err)
	}
	return nil
}

func (r *repositoryDBImpl) Delete(id int) error {
	if err := r.db.Delete(&Model{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}
	return nil
}

func (r *repositoryDBImpl) GetAll() ([]*Instance, error) {
	var models []Model
	if err := r.db.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get all instances: %w", err)
	}

	instances := make([]*Instance, len(models))
	for i, model := range models {
		instances[i] = model.ToEntity()
	}
	return instances, nil
}

func (r *repositoryDBImpl) FindByCategory(category string) ([]*Instance, error) {
	var models []Model
	if err := r.db.Where("category = ?", category).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find instances by category: %w", err)
	}

	instances := make([]*Instance, len(models))
	for i, model := range models {
		instances[i] = model.ToEntity()
	}
	return instances, nil
}
