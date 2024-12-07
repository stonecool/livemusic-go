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

func toInstance(model *types.Model) types.Chrome {
	return &instance.Instance{
		IP:          model.IP,
		Port:        model.Port,
		DebuggerURL: model.DebuggerURL,
		State:       types.InstanceState(model.State),
	}
}

func createInstance(ip string, port int, debuggerURL string, state types.InstanceState) (types.Chrome, error) {
	newInstance := &instance.Instance{
		IP:          ip,
		Port:        port,
		DebuggerURL: debuggerURL,
		State:       state,
	}

	if err := newInstance.Initialize(); err != nil {
		return nil, err
	}

	return newInstance, nil
}

func GetAll() ([]types.Chrome, error) {
	models, err := storage.Repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all instances: %w", err)
	}

	chromes := make([]types.Chrome, len(models))
	for i, m := range models {
		chromes[i] = toInstance(m)
	}

	return chromes, nil
}

func ExistsByIPAndPort(ip string, port int) (bool, error) {
	return storage.Repo.ExistsByIPAndPort(ip, port)
}

// Create Create a local chrome instance
func Create() (types.Chrome, error) {
	ip := "127.0.0.1"

	port, err := util.FindAvailablePort(9222)
	if err != nil {
		internal.Logger.Error("failed to find available port",
			zap.Error(err))
		return nil, err
	}

	if !util.IsValidIPv4(ip) || !util.IsValidPort(port) {
		return nil, fmt.Errorf("invalid ip or port")
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

	internal.Logger.Info("using port for new instance",
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

	chrome, err := createInstance(ip, port, url, types.InstanceStateAvailable)
	if err != nil {
		return nil, err
	}

	if err := chrome.Initialize(); err != nil {
		return nil, err
	}

	return chrome, nil
}

// Bind binds to an existing chrome instance
func Bind(ip string, port int) (types.Chrome, error) {
	if ip == "localhost" {
		ip = "127.0.0.1"
	}

	if !util.IsValidIPv4(ip) || !util.IsValidPort(port) {
		return nil, fmt.Errorf("invalid ip or port")
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

	model, err := storage.Repo.Create(ip, port, url, types.InstanceStateAvailable)
	if err != nil {
		return nil, err
	}

	chrome := toInstance(model)
	if err := chrome.Initialize(); err != nil {
		return nil, err
	}

	return chrome, nil
}
