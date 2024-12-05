package storage

import (
	"fmt"
	"sync"

	"github.com/stonecool/livemusic-go/internal/chrome/types"
)

type MockRepository struct {
	mu     sync.RWMutex
	models map[int]*types.Model
	nextID int
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		models: make(map[int]*types.Model),
		nextID: 1,
	}
}

func (r *MockRepository) Create(ip string, port int, debuggerURL string, state types.ChromeState) (*types.Model, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	model := &types.Model{
		IP:          ip,
		Port:        port,
		DebuggerURL: debuggerURL,
		State:       int(state),
	}

	if err := model.Validate(); err != nil {
		return nil, err
	}

	model.ID = r.nextID
	r.nextID++
	r.models[model.ID] = model

	return model, nil
}

func (r *MockRepository) Get(id int) (*types.Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if model, ok := r.models[id]; ok {
		return model, nil
	}
	return nil, fmt.Errorf("model not found: %d", id)
}

func (r *MockRepository) Update(model *types.Model) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := model.Validate(); err != nil {
		return err
	}

	if _, ok := r.models[model.ID]; !ok {
		return fmt.Errorf("model not found: %d", model.ID)
	}

	r.models[model.ID] = model
	return nil
}

func (r *MockRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.models[id]; !ok {
		return fmt.Errorf("model not found: %d", id)
	}

	delete(r.models, id)
	return nil
}

func (r *MockRepository) GetAll() ([]*types.Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	models := make([]*types.Model, 0, len(r.models))
	for _, model := range r.models {
		models = append(models, model)
	}
	return models, nil
}

func (r *MockRepository) ExistsByIPAndPort(ip string, port int) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, model := range r.models {
		if model.IP == ip && model.Port == port {
			return true, nil
		}
	}
	return false, nil
}
