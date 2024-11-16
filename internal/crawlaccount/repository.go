package crawlaccount

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

func (r *repositoryDBImpl) Create(account *CrawlAccount) error {
	model := &model{}
	model.fromEntity(account)

	if err := r.db.Create(model).Error; err != nil {
		return fmt.Errorf("failed to create crawlaccount: %w", err)
	}

	account.ID = model.ID
	return nil
}

func (r *repositoryDBImpl) Get(id int) (*CrawlAccount, error) {
	var model model
	if err := r.db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("crawlaccount not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get crawlaccount: %w", err)
	}

	return model.toEntity(), nil
}

func (r *repositoryDBImpl) Update(account *CrawlAccount) error {
	model := &model{}
	model.fromEntity(account)

	if err := r.db.Save(model).Error; err != nil {
		return fmt.Errorf("failed to update crawlaccount: %w", err)
	}
	return nil
}

func (r *repositoryDBImpl) Delete(id int) error {
	if err := r.db.Delete(&model{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete crawlaccount: %w", err)
	}
	return nil
}

func (r *repositoryDBImpl) GetAll() ([]*CrawlAccount, error) {
	var models []model
	if err := r.db.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to get all accounts: %w", err)
	}

	accounts := make([]*CrawlAccount, len(models))
	for i, model := range models {
		accounts[i] = model.toEntity()
	}
	return accounts, nil
}

func (r *repositoryDBImpl) FindByCategory(category string) ([]*CrawlAccount, error) {
	var models []model
	if err := r.db.Where("category = ?", category).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find accounts by category: %w", err)
	}

	accounts := make([]*CrawlAccount, len(models))
	for i, model := range models {
		accounts[i] = model.toEntity()
	}
	return accounts, nil
}

func (r *repositoryDBImpl) FindByInstance(instanceID int) ([]*CrawlAccount, error) {
	var models []model
	if err := r.db.Where("instance_id = ?", instanceID).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find accounts by instance: %w", err)
	}

	accounts := make([]*CrawlAccount, len(models))
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
