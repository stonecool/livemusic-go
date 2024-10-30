package internal

var dataType2StructMap = map[string]IDataTable{
	"livehouse": &Livehouse{},
}

type AccountState int

const (
	AS_EXPIRED AccountState = iota // 已过期
	AS_RUNNING                     // 正在运行
)

type InstanceStatus uint8

const (
	INS_OK InstanceStatus = iota
	INS_ERR
)
