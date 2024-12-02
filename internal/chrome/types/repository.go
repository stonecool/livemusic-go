package types

type Repository interface {
	Create(ip string, port int, debuggerURL string, state ChromeState) (*Model, error)
	Get(int) (*Model, error)
	Update(*Model) error
	Delete(int) error
	GetAll() ([]*Model, error)
	ExistsByIPAndPort(ip string, port int) (bool, error)
}
