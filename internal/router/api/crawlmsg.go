package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal"
	http2 "github.com/stonecool/livemusic-go/internal/http"
	"net/http"
)

type addCrawlMsg struct {
	DataType  string `json:"data_type" valid:"Required;MaxSize(100)"`
	DataId    int    `json:"data_id" valid:"Required"`
	CrawlType string `json:"crawl_type" valid:"Required;MaxSize(100)"`
	AccountId string `json:"account_id" valid:"Required;MaxSize(100)"`
}

// AddCrawlMsg
func AddCrawlMsg(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    addCrawlMsg
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if account, err := internal.AddCrawlMsg(form.DataType, form.DataId, form.CrawlType, form.AccountId); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, account)
	}
}

// DeleteCrawlMsg
func DeleteCrawlMsg(ctx *gin.Context) {

}

// ModifyCrawlMsg
func ModifyCrawlMsg(ctx *gin.Context) {

}

// StartCrawlMsgProducer
func StartCrawlMsgProducer(ctx *gin.Context) {

}
