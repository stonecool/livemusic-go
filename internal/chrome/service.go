package chrome

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/database"
	"go.uber.org/zap"
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

func createChrome(ip string, port int, debuggerURL string, state chromeState) (*Chrome, error) {
	return repo.create(ip, port, debuggerURL, state)
}

func getChrome(id int) (*Chrome, error) {
	return repo.get(id)
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
	if err := chrome.RetryInitialize(3); err != nil {
		internal.Logger.Error("failed to reinitialize zombie chrome",
			zap.Error(err),
			zap.String("addr", chrome.getAddr()))
	}

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

	chrome, err := createChrome(ip, port, url, chromeStateConnected)
	if err != nil {
		return nil, err
	}

	if err := chrome.RetryInitialize(3); err != nil {
		internal.Logger.Error("failed to reinitialize zombie chrome",
			zap.Error(err),
			zap.String("addr", chrome.getAddr()))
	}

	err = globalPool.AddChrome(chrome)
	if err != nil {
		return nil, err
	}
	return chrome, nil
}
