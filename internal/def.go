package internal

var dataType2StructMap = map[string]IDataTable{
	"livehouse": &Livehouse{},
}

type AccountState int

const (
	AS_EXPIRED AccountState = iota // 已过期
	AS_RUNNING                     // 正在运行
)

type InstanceState uint8

const (
	INSS_OK InstanceState = iota
	INSS_ERR
)
