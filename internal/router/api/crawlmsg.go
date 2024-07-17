package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal"
	http2 "github.com/stonecool/livemusic-go/internal/http"
	"github.com/unknwon/com"
	"net/http"
)

type crawlMsgForm struct {
	DataType  string `json:"data_type" valid:"Required;MaxSize(100)"`
	DataId    int    `json:"data_id" valid:"Required"`
	CrawlType string `json:"crawl_type" valid:"Required;MaxSize(100)"`
	AccountId string `json:"account_id" valid:"Required;MaxSize(100)"`
}

// AddCrawlMsg
// @Summary	Adds crawl message
// @Accept		json
// @Param		form	body	api.crawlMsgForm	true	"created crawl message"
// @Produce	json
// @Success	200	{object}			http.Response
// @Failure	400	{object}			http.Response
// @Router		/api/v1/crawl-messages 	[post]
func AddCrawlMsg(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    crawlMsgForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if msgProducer, err := internal.AddCrawlMsg(form.DataType, form.DataId, form.CrawlType, form.AccountId); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, msgProducer)
	}
}

// DeleteCrawlMsg
// @Summary	Delete crawl message
// @Param		id	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/crawl-messages/{id} [delete]
func DeleteCrawlMsg(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    idForm
	)

	form.ID = com.StrTo(ctx.Param("id")).MustInt()
	httpCode, errCode := Valid(&form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if internal.DeleteCrawlMsg(form.ID) {
		context.Response(http.StatusOK, 0, nil)
		return
	} else {
		context.Response(http.StatusBadRequest, 0, nil)
	}
}

// ModifyCrawlMsg
// @Summary	Modify crawl message
// @Accept		json
// @Param		form	body	api.crawlMsgForm	true	"modify crawl message"
// @Param		id	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}			http.Response
// @Failure	400	{object}			http.Response
// @Router		/api/v1/crawl-messages/{id} 	[put]
func ModifyCrawlMsg(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    crawlMsgForm
		idForm  idForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	idForm.ID = com.StrTo(ctx.Param("id")).MustInt()
	httpCode, errCode = Valid(&idForm)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if msgProducer, err := internal.AddCrawlMsg(form.DataType, form.DataId, form.CrawlType, form.AccountId); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, msgProducer)
	}
}

// StartCrawlMsgProducer
func StartCrawlMsgProducer(ctx *gin.Context) {

}
