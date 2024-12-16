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
	"time"

	"github.com/stonecool/livemusic-go/internal"
	"go.uber.org/zap"
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
	for port := startPort; port < 65535; port++ {
		if checkPortAvailable(port) {
			return port, nil
		}
	}

	return 0, fmt.Errorf("no available port found")
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
		internal.Logger.Error("failed to start Google Chrome",
			zap.Error(err),
			zap.Int("port", port))
		return err
	}

	time.Sleep(5)
	internal.Logger.Info("Google Chrome started successfully",
		zap.Int("port", port))

	return nil
}

func checkChromeHealth(addr string) (bool, string) {
	url := fmt.Sprintf("http://%s/json/version", addr)
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

func IsValidIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() != nil
}

func IsValidPort(port int) bool {
	return port >= 0 && port <= 65535
}

func GetAddr(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}
