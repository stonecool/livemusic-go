package storage

import (
	"fmt"

	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"github.com/stonecool/livemusic-go/internal/database"
	"gorm.io/gorm"
)

var (
	Repo types.Repository
)

func init() {
	Repo = newRepositoryDB(database.DB)
}

type repository struct {
	db database.Repository[*types.Model]
}

func newRepositoryDB(db *gorm.DB) types.Repository {
	return &repository{
		db: database.NewBaseRepository[*types.Model](db),
	}
}

func (r *repository) Create(ip string, port int, debuggerURL string, state types.ChromeState) (*types.Model, error) {
	model := &types.Model{
		IP:          ip,
		Port:        port,
		DebuggerURL: debuggerURL,
		State:       int(state),
	}

	if err := model.Validate(); err != nil {
		return nil, err
	}

	if err := r.db.Create(model); err != nil {
		return nil, fmt.Errorf("failed to create instance: %w", err)
	}

	return model, nil
}

func (r *repository) Get(id int) (*types.Model, error) {
	if m, err := r.db.Get(id); err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func (r *repository) Update(model *types.Model) error {
	if err := model.Validate(); err != nil {
		return err
	}

	if err := r.db.Update(model).Error; err != nil {
		return fmt.Errorf("failed to update instance: %w", err)
	}
	return nil
}

func (r *repository) Delete(id int) error {
	if err := r.db.Delete(id).Error; err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}
	return nil
}

func (r *repository) GetAll() ([]*types.Model, error) {
	models, err := r.db.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all instances: %w", err)
	}

	return models, nil
}

func (r *repository) ExistsByIPAndPort(ip string, port int) (bool, error) {
	return r.db.ExistsBy("ip = ? AND port = ?", ip, port)
}
