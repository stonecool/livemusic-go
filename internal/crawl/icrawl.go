package crawl

type ICrawl interface {
	GetId() string

	GetName() string

	//GetChan() chan CmdRequest
	//
	//GetState() State
	//
	//SetState(state State)

	Login() (bool, error)

	CheckLogin() (bool, error)

	//Crawl(instance *Instance) error
	//
	//GetLoginRequestData() []byte
	//
	//LoginRequestCallback(request *colly.Request) error
	//
	//LoginResponseCallback(response *colly.Response) error
}
