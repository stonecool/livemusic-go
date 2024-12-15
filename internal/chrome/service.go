package chrome

import (
	"fmt"
	"time"

	"github.com/stonecool/livemusic-go/internal/chrome/pool"

	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/chrome/instance"
	"github.com/stonecool/livemusic-go/internal/chrome/storage"
	"github.com/stonecool/livemusic-go/internal/chrome/types"
	"github.com/stonecool/livemusic-go/internal/chrome/util"
	"go.uber.org/zap"
)

func newInstance(model *types.Model, instanceType types.InstanceType, init bool) types.Chrome {
	ins := &instance.Instance{
		IP:          model.IP,
		Port:        model.Port,
		DebuggerURL: model.DebuggerURL,
		State:       types.InstanceState(model.State),
		Type:        instanceType,
		Opts: &types.InstanceOptions{
			InitTimeout:       time.Second,
			HeartbeatInterval: time.Second * 5,
		},
		StateChan: make(chan types.StateEvent, 1),
	}

	if init {
		ins.Initialize()

		if err := pool.GetPool().AddChrome(ins); err != nil {
			ins.Close()
		}
	}

	return ins
}

func GetAll() []types.Chrome {
	models, err := storage.Repo.GetAll()
	if err != nil {
		internal.Logger.Error("failed to get all instances",
			zap.Error(err))
		return nil
	}

	chromes := make(map[string]types.Chrome)
	for _, m := range models {
		ins := newInstance(m, types.InstanceTypePersistent, false)
		chromes[ins.GetAddr()] = ins
	}

	for _, ins := range pool.GetPool().GetAllChromes() {
		if _, ok := chromes[ins.GetAddr()]; !ok {
			chromes[ins.GetAddr()] = ins
		}
	}

	list := make([]types.Chrome, 0, len(chromes))
	for _, ins := range chromes {
		list = append(list, ins)
	}

	return list
}

func ExistsByIPAndPort(ip string, port int) (bool, error) {
	return storage.Repo.ExistsByIPAndPort(ip, port)
}

// Create Create a local chrome instance
func Create() (types.Chrome, error) {
	port, err := util.FindAvailablePort(9222)
	if err != nil {
		internal.Logger.Error("failed to find available port",
			zap.Error(err))
		return nil, err
	}

	if !util.IsValidPort(port) {
		return nil, fmt.Errorf("invalid ip or port")
	}

	ip := "127.0.0.1"
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

	addr := util.GetAddr(ip, port)
	if pool.GlobalPool.GetChrome(addr) != nil {
		return nil, fmt.Errorf("instance:%s in pool", addr)
	}

	internal.Logger.Info("using port for new instance",
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

	mockModel := &types.Model{
		IP:          ip,
		Port:        port,
		DebuggerURL: url,
		State:       int(types.InstanceStateAvailable),
	}
	chrome := newInstance(mockModel, types.InstanceTypeTemporary, true)

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

	addr := util.GetAddr(ip, port)
	if pool.GlobalPool.GetChrome(addr) != nil {
		return nil, fmt.Errorf("instance:%s in pool", addr)
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

	chrome := newInstance(model, types.InstanceTypePersistent, true)

	return chrome, nil
}
