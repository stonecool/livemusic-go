package instance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/cache"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var instanceCache *cache.Memo

func init() {
	instanceCache = cache.New(getInstance)
}

func getInstance(id int) (interface{}, error) {
	return nil, nil
}

func getInstance1(ip string, port int) (interface{}, error) {
	return nil, nil
}

// Instance
type Instance struct {
	Id          int
	IP          string
	Port        int
	accounts    map[string]*internal.CrawlAccount
	DebuggerUrl string
	Status      internal.InstanceStatus
}

func InitInstance(m *instanceModel) *Instance {
	return &Instance{
		Id:          m.ID,
		IP:          m.IP,
		Port:        m.Port,
		accounts:    map[string]*internal.CrawlAccount{},
		DebuggerUrl: m.DebuggerUrl,
		Status:      internal.InstanceStatus(m.Status),
	}
}

func (i *Instance) getAccounts() map[string]*internal.CrawlAccount {
	return i.accounts
}

func (i *Instance) isAvailable(cat string) bool {
	account, exists := i.accounts[cat]
	if !exists {
		return false
	}

	return account.IsAvailable()
}

func (i *Instance) getAddr() string {
	return i.IP + strconv.Itoa(i.Port)
}

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

// RetryCheckChromeHealth 带重试机制的健康检查
func RetryCheckChromeHealth(addr string, retryCount int, retryDelay time.Duration) (bool, string) {
	for i := 0; i < retryCount; i++ {
		if ok, url := checkChromeHealth(addr); ok {
			return true, url
		}
		time.Sleep(retryDelay)
	}
	return false, ""
}
