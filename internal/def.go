package internal

var DataType2StructMap = map[string]IDataTable{
	"livehouse": &Livehouse{},
}

type InstanceStatus uint

const (
	STATUS_INIT         InstanceStatus = iota // 已创建但未启动
	STATUS_CONNECTED                          // 已连接（启动成功）
	STATUS_DISCONNECTED                       // 已断开
)

type AccountStatus uint8

const (
	ACC_OK AccountStatus = iota
)

type InstanceState uint8

const (
	STATE_UNINITIALIZED InstanceState = iota // 未初始化：实例刚创建
	STATE_INIT_FAILED                        // 初始化失败：start 方法执行失败
	STATE_CONNECTED                          // 连接成功：包含初始化成功和心跳检查正常
	STATE_DISCONNECTED                       // 连接断开：心跳检查失败
)

// 实例事件
type InstanceEvent uint8

const (
	EVENT_START                InstanceEvent = iota // 开始初始化
	EVENT_INIT_SUCCESS                              // 初始化成功
	EVENT_INIT_FAIL                                 // 初始化失败
	EVENT_HEALTH_CHECK_SUCCESS                      // 心跳检查成功
	EVENT_HEALTH_CHECK_FAIL                         // 心跳检查失败
	EVENT_GET_STATE                                 // 获取状态
)
