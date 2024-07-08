package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal/crawl"
	http2 "github.com/stonecool/livemusic-go/internal/http"
	"net/http"
)

type addCoroutineForm struct {
	DataType  string `json:"data_type" valid:"Required;MaxSize(100)"`
	DataId    int    `json:"data_id" valid:"Required"`
	CrawlType string `json:"crawl_type" valid:"Required;MaxSize(100)"`
	AccountId string `json:"account_id" valid:"Required;MaxSize(100)"`
}

// AddCoroutine
func AddCoroutine(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    addCoroutineForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if account, err := crawl.AddCoroutine(form.DataType, form.DataId, form.CrawlType, form.AccountId); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, account)
	}
}

// DeleteCoroutine
func DeleteCoroutine(ctx *gin.Context) {

}

// DeleteCoroutine
func ModifyCoroutine(ctx *gin.Context) {

}

// StartCoroutine
func StartCoroutine(ctx *gin.Context) {

}
