package internal

type Type uint8

const (
	WxPublicAccount Type = iota // wei xin public account
)

var crawlTemplateTypeMap map[Type]string

func init() {
	crawlTemplateTypeMap = make(map[Type]string)

	crawlTemplateTypeMap[WxPublicAccount] = "https://mp.weixin.qq.com/cgi-bin/appmsg"
}

type AccountTemplateType uint8

const (
	NoLogin AccountTemplateType = iota // no need login before Crawl
	Login                              // need login before Crawl
)

type CrawlState uint8

const (
	CsNotLoggedIn CrawlState = iota
	CsLogged
	CsLoggedExpired
)

type Cmd uint8

const (
	CmdReady   Cmd = iota //
	CmdRun                // run a server
	CmdSuspend            // suspend a running server
	CmdCrawl
)

type CmdRequest struct {
	cmd Cmd
}

type State uint8

const (
	StateInitial State = iota
	StateReady
	StateRunning
	StateBlocked
)
