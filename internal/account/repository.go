package account

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/database"
	"gorm.io/gorm"
)

var repo database.Repository[*accountModel]

func init()  {
	repo = newRepositoryDB(database.DB)
}

type repositoryDB struct {
	*database.BaseRepository[*accountModel]
}

func newRepositoryDB(db *gorm.DB) *repositoryDB {
	return &repositoryDB{
		BaseRepository: database.NewBaseRepository[*accountModel](db),
	}
}

func GetAccountByID(id int) (*Account, error) {
	m, err := repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return m.toEntity(), nil
}

func CreateAccountInDB(category string) (*Account, error) {
	m := &accountModel{Category: category}
	if err := m.Validate(); err != nil {
		return nil, err
	}

	if err := repo.Create(m); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return m.toEntity(), nil
}