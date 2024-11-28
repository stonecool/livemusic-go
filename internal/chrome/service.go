package chrome

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/chrome/instance"
	"github.com/stonecool/livemusic-go/internal/chrome/pool"
	"github.com/stonecool/livemusic-go/internal/chrome/storage"
	"github.com/stonecool/livemusic-go/internal/chrome/util"
	"go.uber.org/zap"
)

//func GetInstance(id int) (*Chrome, error) {
//	ins, err := chromeCache.Get(id)
//	if err != nil {
//		return nil, err
//	} else {
//		return ins.(*Chrome), nil
//	}
//}

func createChrome(ip string, port int, debuggerURL string, state instance.ChromeState) (*instance.Chrome, error) {
	return storage.Repo.Create(ip, port, debuggerURL, state)
}

func GetChrome(id int) (*instance.Chrome, error) {
	return storage.Repo.Get(id)
}

func UpdateChrome(chrome *instance.Chrome) error {
	return storage.Repo.Update(chrome)
}

func GetAllChrome() ([]*instance.Chrome, error) {
	return storage.Repo.GetAll()
}

func ExistsByIPAndPort(ip string, port int) (bool, error) {
	return storage.Repo.ExistsByIPAndPort(ip, port)
}

// CreateTempChrome Create a local chrome instance
func CreateTempChrome() (*instance.Chrome, error) {
	ip := "127.0.0.1"
	port, err := util.FindAvailablePort(9222)
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
	err = util.StartChromeOnPort(port)
	if err != nil {
		fmt.Printf("Create instance on port:%d error: %v\n", port, err)
		return nil, err
	}

	ok, url := util.RetryCheckChromeHealth(fmt.Sprintf("%s:%d", ip, port), 3, 1)
	if !ok {
		fmt.Printf("Chrome health check error: %v\n", err)
		return nil, err
	}

	chrome := instance.NewChrome(ip, port, url, instance.ChromeStateConnected)
	if err := chrome.RetryInitialize(3); err != nil {
		internal.Logger.Error("failed to reinitialize zombie chrome",
			zap.Error(err),
			zap.String("addr", chrome.GetAddr()))
	}

	err = pool.GlobalPool.AddChrome(chrome)
	if err != nil {
		return nil, err
	}
	return chrome, nil
}

// BindChrome
func BindChrome(ip string, port int) (*instance.Chrome, error) {
	if !util.IsValidIPv4(ip) || !util.IsValidPort(string(port)) {
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

	ok, url := util.RetryCheckChromeHealth(fmt.Sprintf("%s:%d", ip, port), 3, 1)
	if !ok {
		fmt.Printf("Chrome health check error: %v\n", err)
		return nil, err
	}

	chrome, err := createChrome(ip, port, url, instance.ChromeStateConnected)
	if err != nil {
		return nil, err
	}

	if err := chrome.RetryInitialize(3); err != nil {
		internal.Logger.Error("failed to reinitialize zombie chrome",
			zap.Error(err),
			zap.String("addr", chrome.GetAddr()))
	}

	err = pool.GlobalPool.AddChrome(chrome)
	if err != nil {
		return nil, err
	}
	return chrome, nil
}
