package chrome

import (
	"fmt"

	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/chrome/instance"
	"github.com/stonecool/livemusic-go/internal/chrome/pool"
	"github.com/stonecool/livemusic-go/internal/chrome/storage"
	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"github.com/stonecool/livemusic-go/internal/chrome/util"
	"go.uber.org/zap"
)

func createChrome(ip string, port int, debuggerURL string, state types.ChromeState) (types.IChrome, error) {
	return storage.Repo.Create(ip, port, debuggerURL, state)
}

func GetChrome(id int) (types.IChrome, error) {
	return storage.Repo.Get(id)
}

func UpdateChrome(chrome types.IChrome) error {
	return storage.Repo.Update(chrome)
}

func GetAllChrome() ([]types.IChrome, error) {
	return storage.Repo.GetAll()
}

func ExistsByIPAndPort(ip string, port int) (bool, error) {
	return storage.Repo.ExistsByIPAndPort(ip, port)
}

// CreateTempChrome Create a local chrome instance
func CreateTempChrome() (types.IChrome, error) {
	ip := "127.0.0.1"
	port, err := util.FindAvailablePort(9222)
	if err != nil {
		internal.Logger.Error("failed to find available port",
			zap.Error(err))
		return nil, err
	}

	exists, err := ExistsByIPAndPort(ip, port)
	if err != nil {
		internal.Logger.Error("failed to check port existence",
			zap.Error(err),
			zap.String("ip", ip),
			zap.Int("port", port))
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("port:%d occupied", port)
	}

	internal.Logger.Info("using port for new chrome instance",
		zap.String("ip", ip),
		zap.Int("port", port))

	err = util.StartChromeOnPort(port)
	if err != nil {
		internal.Logger.Error("failed to start chrome",
			zap.Error(err),
			zap.Int("port", port))
		return nil, err
	}

	ok, url := util.RetryCheckChromeHealth(fmt.Sprintf("%s:%d", ip, port), 3, 1)
	if !ok {
		internal.Logger.Error("chrome health check failed",
			zap.String("ip", ip),
			zap.Int("port", port))
		return nil, fmt.Errorf("health check failed")
	}

	chrome, err := storage.Repo.Create(ip, port, url, types.ChromeStateConnected)
	if err != nil {
		return nil, err
	}

	// 类型断言以使用 RetryInitialize 方法
	if chromeInstance, ok := chrome.(*instance.Chrome); ok {
		if err := chromeInstance.RetryInitialize(3); err != nil {
			internal.Logger.Error("failed to initialize chrome",
				zap.Error(err),
				zap.String("addr", chromeInstance.GetAddr()))
			return nil, err
		}

		err = pool.GlobalPool.AddChrome(chromeInstance)
		if err != nil {
			return nil, err
		}
	}

	return chrome, nil
}

// BindChrome binds to an existing chrome instance
func BindChrome(ip string, port int) (types.IChrome, error) {
	if !util.IsValidIPv4(ip) || !util.IsValidPort(string(port)) {
		return nil, fmt.Errorf("invalid ip or port")
	}

	if ip == "localhost" {
		ip = "127.0.0.1"
	}

	exists, err := ExistsByIPAndPort(ip, port)
	if err != nil {
		internal.Logger.Error("failed to check port existence",
			zap.Error(err),
			zap.String("ip", ip),
			zap.Int("port", port))
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("port:%d occupied", port)
	}

	ok, url := util.RetryCheckChromeHealth(fmt.Sprintf("%s:%d", ip, port), 3, 1)
	if !ok {
		internal.Logger.Error("chrome health check failed",
			zap.String("ip", ip),
			zap.Int("port", port))
		return nil, fmt.Errorf("health check failed")
	}

	chrome, err := createChrome(ip, port, url, types.ChromeStateConnected)
	if err != nil {
		return nil, err
	}

	// 类型断言以使用 RetryInitialize 方法
	if chromeInstance, ok := chrome.(*instance.Chrome); ok {
		if err := chromeInstance.RetryInitialize(3); err != nil {
			internal.Logger.Error("failed to initialize chrome",
				zap.Error(err),
				zap.String("addr", chromeInstance.GetAddr()))
			return nil, err
		}

		err = pool.GlobalPool.AddChrome(chromeInstance)
		if err != nil {
			return nil, err
		}
	}

	return chrome, nil
}
