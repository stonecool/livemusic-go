package chrome

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stonecool/livemusic-go/internal/model"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func checkPortAvailable(port int) bool {
	var cmd *exec.Cmd
	portStr := strconv.Itoa(port)

	switch runtime.GOOS {
	case "darwin", "linux":
		cmd = exec.Command("lsof", "-i", ":"+portStr)
	case "windows":
		cmd = exec.Command("netstat", "-an", "|", "findstr", ":"+portStr)

	default:
		return false
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		// 如果命令执行失败，假设端口可用
		return true
	}

	return !strings.Contains(out.String(), portStr)
}

func FindAvailablePort(startPort int) (int, error) {
	var wg sync.WaitGroup
	portChan := make(chan int, 1)

	for port := startPort; port < 65535; port++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			if checkPortAvailable(p) {
				select {
				case portChan <- p:
				default:
				}
			}
		}(port)
	}

	go func() {
		wg.Wait()
		close(portChan)
	}()

	select {
	case port := <-portChan:
		return port, nil
	case <-time.After(10 * time.Second):
		return 0, fmt.Errorf("no available port found")
	}
}

func CreateInstance(port int) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		// macOS
		cmd = exec.Command("open", "-a", "Google Chrome", "--args", "--remote-debugging-port="+strconv.Itoa(port))
	case "windows":
		// Windows
		cmd = exec.Command("cmd", "/c", "start", "chrome", "--remote-debugging-port="+strconv.Itoa(port))
	default:
		// Linux
		cmd = exec.Command("google-chrome", "--remote-debugging-port="+strconv.Itoa(port))
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Failed to start Google Chrome: %v\n", err)
		return err
	}

	// 等待几秒钟确保chrome实例启动
	time.Sleep(5)
	fmt.Println("Google Chrome started successfully")

	return nil
}

// checkChromeHealth 通过远程调试端口检查 Chrome 实例的健康状态
func checkChromeHealth(addr string) (bool, string) {
	url := fmt.Sprintf("http://%s/json", addr)
	resp, err := http.Get(url)
	if err != nil {
		return false, ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, ""
	}

	type chromeVersion struct {
		Browser              string `json:"Browser"`
		ProtocolVersion      string `json:"Protocol-Version"`
		UserAgent            string `json:"User-Agent"`
		V8Version            string `json:"V8-Version"`
		WebKitVersion        string `json:"WebKit-Version"`
		WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	}

	var version chromeVersion
	err = json.NewDecoder(resp.Body).Decode(&version)
	if err != nil {
		return false, ""
	}

	debuggerUrl := version.WebSocketDebuggerURL
	return len(debuggerUrl) > 0, debuggerUrl
}

func RetryCheckChromeHealth(addr string, retryCount int, retryDelay time.Duration) (bool, string) {
	for i := 0; i < retryCount; i++ {
		if ok, url := checkChromeHealth(addr); ok {
			return true, url
		}
		time.Sleep(time.Second * retryDelay)
	}
	return false, ""
}

// CreateLocalChromeInstance Create a local chrome instance
func CreateLocalChromeInstance() (*Instance, error) {
	ip := "127.0.0.1"
	port, err := FindAvailablePort(9222)
	if err != nil {
		fmt.Printf("Failed to find an available port: %v\n", err)
		return nil, err
	}

	exists, err := model.ExistsChromeInstance(ip, port)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}

	if exists {
		return nil, fmt.Errorf("port:%d occupied", port)
	}

	fmt.Printf("Using ip:%s port: %d\n", ip, port)
	err = CreateInstance(port)
	if err != nil {
		fmt.Printf("Create instance on port:%d error: %v\n", port, err)
		return nil, err
	}

	ok, url := RetryCheckChromeHealth(fmt.Sprintf("%s:%d", ip, port), 3, 1)
	if !ok {
		fmt.Printf("Instance health check error: %v\n", err)
		return nil, err
	}

	data := map[string]interface{}{
		"ip":           ip,
		"port":         port,
		"debugger_url": url,
	}

	m, err := model.AddChromeInstance(data)
	if err != nil {
		return nil, err
	}

	return globalPool.AddInstance(m.ID)
}

// BindChromeInstance
//func BindChromeInstance(ip string, port int) (*Instance, error) {
//	exists, err := model.ExistsChromeInstance(ip, port)
//	if err != nil {
//		fmt.Printf("%v\n", err)
//		return nil, err
//	}
//
//	if exists {
//		return nil, fmt.Errorf("port:%d occupied", port)
//	}
//
//	ok, _ := RetryCheckChromeHealth(fmt.Sprintf("%s:%d", ip, port), 3, 1)
//	if !ok {
//		fmt.Printf("Instance health check error: %v\n", err)
//		return nil, err
//	}
//
//	return globalPool.AddInstance(m.ID)
//}
