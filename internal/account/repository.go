package account

import (
	"fmt"

	"gorm.io/gorm"
)

type repositoryDBImpl struct {
	db *gorm.DB
}

func NewRepositoryDB(db *gorm.DB) IRepository {
	return &repositoryDBImpl{db: db}
}

func (r *repositoryDBImpl) Create(account *Account) error {
	model := &model{}
	model.fromEntity(account)

	if err := r.db.Create(model).Error; err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	account.ID = model.ID
	return nil
}

func (r *repositoryDBImpl) Get(id int) (*Account, error) {
	var model model
	if err := r.db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("account not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return model.toEntity(), nil
}

func (r *repositoryDBImpl) Update(account *Account) error {
	model := &model{}
	model.fromEntity(account)

	if err := r.db.Save(model).Error; err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}
	return nil
}

func (r *repositoryDBImpl) Delete(id int) error {
	if err := r.db.Delete(&model{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}
	return nil
}

func (r *repositoryDBImpl) GetAll() ([]*Account, error) {
	var models []model
	if err := r.db.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get all accounts: %w", err)
	}

	accounts := make([]*Account, len(models))
	for i, model := range models {
		accounts[i] = model.toEntity()
	}
	return accounts, nil
}

func (r *repositoryDBImpl) FindByCategory(category string) ([]*Account, error) {
	var models []model
	if err := r.db.Where("category = ?", category).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find accounts by category: %w", err)
	}

	accounts := make([]*Account, len(models))
	for i, model := range models {
		accounts[i] = model.toEntity()
	}
	return accounts, nil
}

func (r *repositoryDBImpl) FindByInstance(id int) ([]*Account, error) {
	var models []model
	if err := r.db.Where("instance_id = ?", id).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find accounts by instance: %w", err)
	}

	accounts := make([]*Account, len(models))
	for i, model := range models {
		accounts[i] = model.toEntity()
	}
	return accounts, nil
}

// TODO support transaction
func (r *repositoryDBImpl) WithTx(tx *gorm.DB) IRepository {
	return &repositoryDBImpl{db: tx}
}

func (r *repositoryDBImpl) Transaction(fn func(IRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(r.WithTx(tx))
	})
}
