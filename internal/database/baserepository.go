package database

import (
	"fmt"

	"gorm.io/gorm"
)

type Entity interface {
	GetID() int
	TableName() string
	Validate() error
}

type Repository[T Entity] interface {
	Create(entity T) error
	Get(id int) (T, error)
	Update(entity T) error
	Delete(id int) error
	GetAll() ([]T, error)
	FindBy(query string, args ...interface{}) ([]T, error)
	ExistsBy(query string, args ...interface{}) (bool, error)
}

type BaseRepository[T Entity] struct {
	db *gorm.DB
}

func NewBaseRepository[T Entity](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{
		db: db,
	}
}

func (r *BaseRepository[T]) Create(entity T) error {
	if err := entity.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := r.db.Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create entity: %w", err)
	}

	return nil
}

func (r *BaseRepository[T]) Get(id int) (T, error) {
	var entity T
	if err := r.db.First(&entity, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity, fmt.Errorf("entity not found: %d", id)
		}
		return entity, fmt.Errorf("failed to get entity: %w", err)
	}
	return entity, nil
}

func (r *BaseRepository[T]) Update(entity T) error {
	if err := entity.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := r.db.Save(entity).Error; err != nil {
		return fmt.Errorf("failed to update entity: %w", err)
	}
	return nil
}

func (r *BaseRepository[T]) Delete(id int) error {
	var entity T
	if err := r.db.Delete(&entity, id).Error; err != nil {
		return fmt.Errorf("failed to delete entity: %w", err)
	}
	return nil
}

func (r *BaseRepository[T]) GetAll() ([]T, error) {
	var entities []T
	if err := r.db.Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("failed to find entities: %w", err)
	}
	return entities, nil
}

func (r *BaseRepository[T]) FindBy(query string, args ...interface{}) ([]T, error) {
	var entities []T
	if err := r.db.Where(query, args...).Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("failed to find entities: %w", err)
	}
	return entities, nil
}

func (r *BaseRepository[T]) ExistsBy(query string, args ...interface{}) (bool, error) {
	var count int64
	var entity T
	err := r.db.Model(&entity).Where(query, args...).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
