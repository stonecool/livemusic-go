package chrome

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/database"
)

func CreateChrome(ip string, port int, debuggerURL string) (*Chrome, error) {
	repo := NewRepositoryDB(database.DB)
	factory := NewFactory(repo)
	return factory.CreateChrome(ip, port, debuggerURL)
}

func getChrome(id int) (*Chrome, error) {
	repo := NewRepositoryDB(database.DB)
	chrome, err := repo.Get(id)
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
	repo := NewRepositoryDB(database.DB)
	return repo.GetAll()
}

func ExistsByIPAndPort(ip string, port int) (bool, error) {
	repo := NewRepositoryDB(database.DB)
	return repo.ExistsByIPAndPort(ip, port)
}