package chrome

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/database"
)

var (
	//chromeCache *cache.Memo
	repo repository
)

func init() {
	//chromeCache = cache.New(func(id int) (interface{}, error) {
	//	return getChrome(id)
	//})
	repo = newRepositoryDB(database.DB)
}

//func GetInstance(id int) (*Chrome, error) {
//	ins, err := chromeCache.Get(id)
//	if err != nil {
//		return nil, err
//	} else {
//		return ins.(*Chrome), nil
//	}
//}

func createChrome(ip string, port int, debuggerURL string) (*Chrome, error) {
	return repo.create(ip, port, debuggerURL, chromeStateUninitialized)
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

// CreateTempChrome Create a local chrome instance
func CreateTempChrome() (*Chrome, error) {
	ip := "127.0.0.1"
	port, err := FindAvailablePort(9222)
	if err != nil {
		fmt.Printf("Failed to find an available port: %v\n", err)
		return nil, err
	}

	exists, err := ExistsByIPAndPort(ip, port)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("port:%d occupied", port)
	}

	fmt.Printf("Using ip:%s port: %d\n", ip, port)
	err = startChromeOnPort(port)
	if err != nil {
		fmt.Printf("Create instance on port:%d error: %v\n", port, err)
		return nil, err
	}

	ok, url := RetryCheckChromeHealth(fmt.Sprintf("%s:%d", ip, port), 3, 1)
	if !ok {
		fmt.Printf("Chrome health check error: %v\n", err)
		return nil, err
	}

	chrome := newChrome(ip, port, url, chromeStateConnected)

	//m, err := CreateChrome(ip, port, url)
	//if err != nil {
	//	return nil, err
	//}

	err = globalPool.AddChrome(chrome)
	if err != nil {
		return nil, err
	}
	return chrome, nil
}

// BindChrome
func BindChrome(ip string, port int) (*Chrome, error) {
	if !IsValidIPv4(ip) || !IsValidPort(string(port)) {
		return nil, fmt.Errorf("invalid")
	}

	if ip == "localhost" {
		ip = "127.0.0.1"
	}

	exists, err := ExistsByIPAndPort(ip, port)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("port:%d occupied", port)
	}

	ok, url := RetryCheckChromeHealth(fmt.Sprintf("%s:%d", ip, port), 3, 1)
	if !ok {
		fmt.Printf("Chrome health check error: %v\n", err)
		return nil, err
	}

	chrome, err := createChrome(ip, port, url)
	if err != nil {
		return nil, err
	}

	err = globalPool.AddChrome(chrome)
	if err != nil {
		return nil, err
	}
	return chrome, nil
}
