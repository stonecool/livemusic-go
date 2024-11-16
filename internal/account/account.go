package account

import (
	"sync"

	"github.com/stonecool/livemusic-go/internal"
)

type Account struct {
	ID          int    `json:"id"`
	Category    string `json:"category"`
	AccountName string `json:"account_name"`
	lastURL     string
	cookies     []byte
	InstanceID  int
	State       internal.AccountState
	mu          sync.RWMutex
	msgChan     chan *internal.AsyncMessage // 改用 AsyncMessage
	done        chan struct{}
}

func NewAccount(m *model) *Account {
	acc := &Account{
		model:   m,
		msgChan: make(chan *AsyncMessage),
		done:    make(chan struct{}),
	}
	go acc.processTask()
	return acc
}

// 实现 Account 接口
func (a *Account) GetID() int {
	return a.model.ID
}

func (a *Account) GetCategory() string {
	return a.model.Category
}

func (a *Account) GetState() AccountState {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.model.State
}

func (a *Account) IsAvailable() bool {
	return a.GetState() == AS_RUNNING
}
