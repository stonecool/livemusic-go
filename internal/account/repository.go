package account

import (
	"github.com/stonecool/livemusic-go/internal/database"
	"gorm.io/gorm"
)

type repositoryDB struct {
	*database.BaseRepository[*model]
}

func newRepositoryDB(db *gorm.DB) *repositoryDB {
	return &repositoryDB{
		BaseRepository: database.NewBaseRepository[*model](db),
	}
}