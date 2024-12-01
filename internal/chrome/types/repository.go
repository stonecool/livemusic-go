package types

type Repository interface {
	Create(ip string, port int, debuggerURL string, state ChromeState) (Chrome, error)
	Get(id int) (Chrome, error)
	Update(chrome Chrome) error
	Delete(id int) error
	GetAll() ([]Chrome, error)
	ExistsByIPAndPort(ip string, port int) (bool, error)
}
