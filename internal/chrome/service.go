package chrome

import (
	"fmt"

	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/chrome/instance"
	"github.com/stonecool/livemusic-go/internal/chrome/storage"
	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"github.com/stonecool/livemusic-go/internal/chrome/util"
	"go.uber.org/zap"
)

func createChrome(dto *types.ChromeDTO) (types.Chrome, error) {
	chrome := instance.NewChrome(
		dto.IP,
		dto.Port,
		dto.DebuggerURL,
		dto.State,
	)

	if err := chrome.Initialize(); err != nil {
		return nil, err
	}

	return chrome, nil
}

func GetChrome(id int) (types.Chrome, error) {
	dto, err := storage.Repo.Get(id)
	if err != nil {
		return nil, err
	}

	return createChrome(dto)
}

func UpdateChrome(chrome types.Chrome) error {
	return storage.Repo.Update(chrome)
}

func GetAllChrome() ([]types.Chrome, error) {
	return storage.Repo.GetAll()
}

func ExistsByIPAndPort(ip string, port int) (bool, error) {
	return storage.Repo.ExistsByIPAndPort(ip, port)
}

// CreateTempChrome Create a local chrome instance
func CreateTempChrome() (types.Chrome, error) {
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

	if err := chrome.Initialize(); err != nil {
		return nil, err
	}

	return chrome, nil
}

// BindChrome binds to an existing chrome instance
func BindChrome(ip string, port int) (types.Chrome, error) {
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

	if err := chrome.Initialize(); err != nil {
		return nil, err
	}

	return chrome, nil
}
