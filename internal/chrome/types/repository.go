package types

type ChromeDTO struct {
	ID          int
	IP          string
	Port        int
	DebuggerURL string
	State       ChromeState
}

type Repository interface {
	Create(dto ChromeDTO) (*ChromeDTO, error)
	Get(id int) (*ChromeDTO, error)
	Update(dto ChromeDTO) error
	Delete(id int) error
	GetAll() ([]*ChromeDTO, error)
	ExistsByIPAndPort(ip string, port int) (bool, error)
}
