package account

import (
	"github.com/stonecool/livemusic-go/internal/database"
	"gorm.io/gorm"
)

type AccountRepository1 struct {
	*database.BaseRepository[*model]
}

func NewAccountRepository(db *gorm.DB) *AccountRepository1 {
	return &AccountRepository1{
		BaseRepository: database.NewBaseRepository[*model](db),
	}
}

// 只需要实现特殊的方法
func (r *AccountRepository1) FindByCategory(category string) ([]*Account, error) {
	models, err := r.FindBy("category = ?", category)
	if err != nil {
		return nil, err
	}

	accounts := make([]*Account, len(models))
	for i, model := range models {
		accounts[i] = model.toEntity()
	}
	return accounts, nil
}
