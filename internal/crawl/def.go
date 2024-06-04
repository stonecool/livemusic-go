package crawl

type AccountTemplateType uint8

const (
	NoLogin AccountTemplateType = iota // no need login before Crawl
	Login                              // need login before Crawl
)

type AccountState uint8

const (
	AccountStateUnLogged AccountState = iota
	AccountStateLogged
	AccountStateLoggedExpired
)

type Cmd uint8

const (
	CmdReady   Cmd = iota //
	CmdRun                // run a server
	CmdSuspend            // suspend a running server
	CmdCrawl
)

type CmdRequest struct {
	cmd      Cmd
	instance *Instance
}

type State uint8

const (
	StateInitial State = iota
	StateReady
	StateRunning
	StateBlocked
)
