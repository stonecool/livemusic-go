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

func createChrome(model *types.Model) (types.Chrome, error) {
	chrome := NewInstance(
		model.IP,
		model.Port,
		model.DebuggerURL,
		types.ChromeState(model.State),
	)

	if err := chrome.Initialize(); err != nil {
		return nil, err
	}

	return chrome, nil
}

func createChromeWithParam(ip string, port int, debuggerURL string, state types.ChromeState) (types.Chrome, error) {
	instance := NewInstance(ip, port, debuggerURL, state)
	if err := instance.Initialize(); err != nil {
		return nil, err
	}

	return instance, nil
}

func GetChrome(id int) (types.Chrome, error) {
	model, err := storage.Repo.Get(id)
	if err != nil {
		return nil, err
	}

	return createChrome(model)
}

func GetAllChrome() ([]types.Chrome, error) {
	models, err := storage.Repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all instances: %w", err)

	}

	chromes := make([]types.Chrome, len(models))
	for i, m := range models {
		chromes[i] = modelToChrome(m)
	}

	return chromes, nil
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

	internal.Logger.Info("using port for new model instance",
		zap.String("ip", ip),
		zap.Int("port", port))

	err = util.StartChromeOnPort(port)
	if err != nil {
		internal.Logger.Error("failed to start model",
			zap.Error(err),
			zap.Int("port", port))
		return nil, err
	}

	ok, url := util.RetryCheckChromeHealth(fmt.Sprintf("%s:%d", ip, port), 3, 1)
	if !ok {
		internal.Logger.Error("model health check failed",
			zap.String("ip", ip),
			zap.Int("port", port))
		return nil, fmt.Errorf("health check failed")
	}

	model, err := storage.Repo.Create(ip, port, url, types.ChromeStateConnected)
	if err != nil {
		return nil, err
	}

	ins := modelToChrome(model)
	if err := ins.Initialize(); err != nil {
		return nil, err
	}

	return ins, nil
}

// BindChrome binds to an existing chrome instance
func BindChrome(ip string, port int) (types.Chrome, error) {
	if !util.IsValidIPv4(ip) || !util.IsValidPort(port) {
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
		internal.Logger.Error("model health check failed",
			zap.String("ip", ip),
			zap.Int("port", port))
		return nil, fmt.Errorf("health check failed")
	}

	model, err := createChromeWithParam(ip, port, url, types.ChromeStateConnected)
	if err != nil {
		return nil, err
	}

	if err := model.Initialize(); err != nil {
		return nil, err
	}

	return model, nil
}

func modelToChrome(model *types.Model) types.Chrome {
	return NewInstance(
		model.IP,
		model.Port,
		model.DebuggerURL,
		types.ChromeState(model.State),
	)
}

func NewInstance(ip string, port int, url string, state types.ChromeState) *instance.Instance {
	return &instance.Instance{
		IP:          ip,
		Port:        port,
		DebuggerURL: url,
		State:       state,
	}
}
