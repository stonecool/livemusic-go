package internal

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/chrome"
	"github.com/stonecool/livemusic-go/internal/instance"
)

// CreateLocalChromeInstance Create a local chrome instance
func CreateLocalChromeInstance() (*instance.Instance, error) {
	ip := "127.0.0.1"
	port, err := instance.FindAvailablePort(9222)
	if err != nil {
		fmt.Printf("Failed to find an available port: %v\n", err)
		return nil, err
	}

	exists, err := instance.ExistsChromeInstance(ip, port)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("port:%d occupied", port)
	}

	fmt.Printf("Using ip:%s port: %d\n", ip, port)
	err = instance.CreateInstance(port)
	if err != nil {
		fmt.Printf("Create instance on port:%d error: %v\n", port, err)
		return nil, err
	}

	ok, url := instance.RetryCheckChromeHealth(fmt.Sprintf("%s:%d", ip, port), 3, 1)
	if !ok {
		fmt.Printf("Instance health check error: %v\n", err)
		return nil, err
	}

	data := map[string]interface{}{
		"ip":           ip,
		"port":         port,
		"debugger_url": url,
		"status":       INS_OK,
	}

	m, err := instance.AddChromeInstance(data)
	if err != nil {
		return nil, err
	}

	ins := instance.InitInstance(m)
	chrome.Pool.AddInstance(ins)
	return ins, nil
}

// BindChromeInstance
func BindChromeInstance(ip string, port int) (*instance.Instance, error) {
	exists, err := instance.ExistsChromeInstance(ip, port)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("port:%d occupied", port)
	}

	ok, url := instance.RetryCheckChromeHealth(fmt.Sprintf("%s:%d", ip, port), 3, 1)
	if !ok {
		fmt.Printf("Instance health check error: %v\n", err)
		return nil, err
	}

	data := map[string]interface{}{
		"ip":           ip,
		"port":         port,
		"debugger_url": url,
		"status":       INS_OK,
	}

	m, err := instance.AddChromeInstance(data)
	if err != nil {
		return nil, err
	}

	ins := instance.InitInstance(m)
	chrome.Pool.AddInstance(ins)
	return ins, nil
}
