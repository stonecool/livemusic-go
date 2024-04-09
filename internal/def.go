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
