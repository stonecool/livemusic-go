package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
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

func StartChromeOnPort(port int) error {
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

// ConvertIPv4ToInt converts an IPv4 address to a 32-bit unsigned integer
func ConvertIPv4ToInt(ip string) (uint32, error) {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return 0, fmt.Errorf("invalid IPv4 address")
	}

	var result uint32
	for i := 0; i < 4; i++ {
		part, err := strconv.Atoi(parts[i])
		if err != nil {
			return 0, err
		}
		result = result<<8 + uint32(part)
	}
	return result, nil
}

// CombineIPAndPort combines an IP address and port into a unique 64-bit unsigned integer
func CombineIPAndPort(ip string, port uint16) (uint64, error) {
	ipInt, err := ConvertIPv4ToInt(ip)
	if err != nil {
		return 0, err
	}
	return (uint64(ipInt) << 16) | uint64(port), nil
}

func IsValidIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() != nil
}

func IsValidPort(port string) bool {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	return portNum >= 0 && portNum <= 65535
}