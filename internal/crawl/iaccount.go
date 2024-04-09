package crawl

import "github.com/gocolly/colly"

type IAccount interface {
	GetId() int

	GetChan() chan CmdRequest

	GetState() State

	SetState(state State)

	Login() (bool, error)

	Crawl(instance *Instance) error

	GetLoginRequestData() []byte

	LoginRequestCallback(request *colly.Request) error

	LoginResponseCallback(response *colly.Response) error
}
