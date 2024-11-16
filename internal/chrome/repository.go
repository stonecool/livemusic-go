package chrome

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	Create(chrome *Chrome) error
	Get(id int) (*Chrome, error)
	Update(chrome *Chrome) error
	Delete(id int) error
	GetAll() ([]*Chrome, error)
}

type repositoryDBImpl struct {
	db *gorm.DB
}

func NewRepositoryDB(db *gorm.DB) Repository {
	return &repositoryDBImpl{db: db}
}

func (r *repositoryDBImpl) Create(chrome *Chrome) error {
	m := &model{}
	m.fromEntity(chrome)

	if err := r.db.Create(m).Error; err != nil {
		return fmt.Errorf("failed to create instance: %w", err)
	}

	chrome.ID = m.ID
	return nil
}

func (r *repositoryDBImpl) Get(id int) (*Chrome, error) {
	var m model
	if err := r.db.First(&m, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("instance not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}

	return m.toEntity(), nil
}

func (r *repositoryDBImpl) Update(chrome *Chrome) error {
	m := &model{}
	m.fromEntity(chrome)

	if err := r.db.Save(m).Error; err != nil {
		return fmt.Errorf("failed to update instance: %w", err)
	}
	return nil
}

func (r *repositoryDBImpl) Delete(id int) error {
	if err := r.db.Delete(&model{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}
	return nil
}

func (r *repositoryDBImpl) GetAll() ([]*Chrome, error) {
	var models []*model
	if err := r.db.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get all chromes: %w", err)
	}

	chromes := make([]*Chrome, len(models))
	for i, m := range models {
		chromes[i] = m.toEntity()
	}

	return chromes, nil
}
