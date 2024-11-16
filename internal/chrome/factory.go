package chrome

import (
	"fmt"

	"gorm.io/gorm"
)

type Factory interface {
	CreateInstance(category string) (*Instance, error)
}

type factoryImpl struct {
	repo Repository
}

func NewFactory(repo Repository) Factory {
	return &factoryImpl{repo: repo}
}

func (f *factoryImpl) CreateInstance(category string) (*Instance, error) {
	v := NewValidator()
	if err := v.ValidateCategory(category); err != nil {
		return nil, fmt.Errorf("invalid instance category: %w", err)
	}

	instance, err := NewChromeInstance(category)
	if err != nil {
		return nil, fmt.Errorf("failed to create chrome instance: %w", err)
	}

	if err := f.repo.Create(instance); err != nil {
		return nil, fmt.Errorf("failed to save instance: %w", err)
	}

	return instance, nil
}

// 便捷创建方法
func CreateInstance1(db *gorm.DB, category string) (*Instance, error) {
	repo := NewRepositoryDB(db)
	factory := NewFactory(repo)
	return factory.CreateInstance(category)
}
