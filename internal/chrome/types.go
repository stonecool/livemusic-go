package chrome

// ChromeState 表示 Chrome 实例的状态
type ChromeState uint8

const (
	ChromeStateConnected    ChromeState = iota // 连接成功：包含初始化成功和心跳检查正常
	ChromeStateDisconnected                    // 连接断开：心跳检查失败
	ChromeStateOffline
)
