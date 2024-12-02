package storage

import (
	"fmt"

	"github.com/stonecool/livemusic-go/internal/chrome/types"

	"github.com/stonecool/livemusic-go/internal/chrome/instance"
	"github.com/stonecool/livemusic-go/internal/database"
	"gorm.io/gorm"
)

var (
	Repo types.Repository
)

func init() {
	Repo = newRepositoryDB(database.DB)
}

type repositoryDB struct {
	db database.Repository[*model]
}

func newRepositoryDB(db *gorm.DB) types.Repository {
	return &repositoryDB{
		db: database.NewBaseRepository[*model](db),
	}
}

func (r *repositoryDB) Create(dto types.ChromeDTO) (*types.ChromeDTO, error) {
	m := &model{
		IP:          dto.IP,
		Port:        dto.Port,
		DebuggerURL: dto.DebuggerURL,
		State:       int(dto.State),
	}

	if err := r.db.Create(m); err != nil {
		return nil, fmt.Errorf("failed to create chrome: %w", err)
	}

	return m.toDTO(), nil
}

func (r *repositoryDB) Get(id int) (*types.ChromeDTO, error) {
	m, err := r.db.Get(id)
	if err != nil {
		return nil, err
	}
	return m.toDTO(), nil
}

func (r *repositoryDB) Update(chrome types.Chrome) error {
	// 由于接口转换的问题，这里需要做一个类型断言
	chromeInstance, ok := chrome.(*instance.Chrome)
	if !ok {
		return fmt.Errorf("invalid chrome instance type")
	}

	m := &model{}
	m.fromEntity(chromeInstance)

	if err := r.db.Update(m).Error; err != nil {
		return fmt.Errorf("failed to update instance: %w", err)
	}
	return nil
}

func (r *repositoryDB) Delete(id int) error {
	if err := r.db.Delete(id).Error; err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}
	return nil
}

func (r *repositoryDB) GetAll() ([]types.Chrome, error) {
	models, err := r.db.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all instances: %w", err)
	}

	chromes := make([]types.Chrome, len(models))
	for i, m := range models {
		chromes[i] = m.toEntity()
	}

	return chromes, nil
}

func (r *repositoryDB) ExistsByIPAndPort(ip string, port int) (bool, error) {
	return r.db.ExistsBy("ip = ? AND port = ?", ip, port)
}
