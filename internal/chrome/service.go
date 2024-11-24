package chrome

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/cache"
	"github.com/stonecool/livemusic-go/internal/database"
)

var (
	chromeCache *cache.Memo
	repo        repository
)

func init() {
	chromeCache = cache.New(func(id int) (interface{}, error) {
		return getChrome(id)
	})
	repo = newRepositoryDB(database.DB)
}

func GetInstance(id int) (*Chrome, error) {
	ins, err := chromeCache.Get(id)
	if err != nil {
		return nil, err
	} else {
		return ins.(*Chrome), nil
	}
}

func CreateChrome(ip string, port int, debuggerURL string) (*Chrome, error) {
	return repo.create(ip, port, debuggerURL, STATE_UNINITIALIZED)
}

func getChrome(id int) (*Chrome, error) {
	chrome, err := repo.get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get chrome instance: %w", err)
	}

	if chrome.NeedsReInitialize() {
		if err := chrome.RetryInitialize(3); err != nil {
			return nil, fmt.Errorf("failed to reinitialize chrome instance: %w", err)
		}
	}

	return chrome, nil
}

func GetAllChrome() ([]*Chrome, error) {
	return repo.getAll()
}

func ExistsByIPAndPort(ip string, port int) (bool, error) {
	return repo.existsByIPAndPort(ip, port)
}
