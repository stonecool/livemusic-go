package account

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	Create(account *Account) error
	Get(id int) (*Account, error)
	Update(account *Account) error
	Delete(id int) error
	GetAll() ([]*Account, error)
	FindByCategory(category string) ([]*Account, error)
	FindByInstance(instanceID int) ([]*Account, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(account *Account) error {
	model := &Model{}
	model.FromEntity(account)

	if err := r.db.Create(model).Error; err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	account.ID = model.ID
	return nil
}

func (r *repository) Get(id int) (*Account, error) {
	var model Model
	if err := r.db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("account not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return model.ToEntity(), nil
}

func (r *repository) Update(account *Account) error {
	model := &Model{}
	model.FromEntity(account)

	if err := r.db.Save(model).Error; err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}
	return nil
}

func (r *repository) Delete(id int) error {
	if err := r.db.Delete(&Model{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}
	return nil
}

func (r *repository) GetAll() ([]*Account, error) {
	var models []Model
	if err := r.db.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get all accounts: %w", err)
	}

	accounts := make([]*Account, len(models))
	for i, model := range models {
		accounts[i] = model.ToEntity()
	}
	return accounts, nil
}

func (r *repository) FindByCategory(category string) ([]*Account, error) {
	var models []Model
	if err := r.db.Where("category = ?", category).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find accounts by category: %w", err)
	}

	accounts := make([]*Account, len(models))
	for i, model := range models {
		accounts[i] = model.ToEntity()
	}
	return accounts, nil
}

func (r *repository) FindByInstance(instanceID int) ([]*Account, error) {
	var models []Model
	if err := r.db.Where("instance_id = ?", instanceID).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find accounts by instance: %w", err)
	}

	accounts := make([]*Account, len(models))
	for i, model := range models {
		accounts[i] = model.ToEntity()
	}
	return accounts, nil
}
