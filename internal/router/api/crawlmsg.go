package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal"
	http2 "github.com/stonecool/livemusic-go/internal/http"
	"github.com/unknwon/com"
	"net/http"
)

type crawlMsgForm struct {
	DataType        string `json:"data_type" valid:"Required;MaxSize(100)"`
	DataId          int    `json:"data_id" valid:"Required"`
	AccountType     string `json:"account_type" valid:"Required;MaxSize(100)"`
	TargetAccountId string `json:"target_account_id" valid:"Required;MaxSize(100)"`
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

	msg := internal.CrawlMsg{
		DataType:        form.DataType,
		DataId:          form.DataId,
		AccountType:     form.AccountType,
		TargetAccountId: form.TargetAccountId,
	}

	if err := msg.Add(); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, msg)
	}
}

// GetCrawlMsg
// @Summary	Get a crawl message
// @Param		id	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/crawl-messages/{ID} [get]
func GetCrawlMsg(ctx *gin.Context) {
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

	msg := internal.CrawlMsg{ID: form.ID}
	if err := msg.Get(); err != nil {
		context.Response(http.StatusBadRequest, 0, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, msg)
	}
}

// GetCrawlMsgs
// @Summary	Get all crawl messages
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/crawl-messages [get]
func GetCrawlMsgs(ctx *gin.Context) {
	var context = http2.Context{Context: ctx}

	msg := internal.CrawlMsg{}
	if msgs, err := msg.GetAll(); err != nil {
		context.Response(http.StatusBadRequest, 0, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, msgs)
	}
}

// DeleteCrawlMsg
// @Summary	Delete crawl message
// @Param		id	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/crawl-messages/{ID} [delete]
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

	msg := &internal.CrawlMsg{ID: form.ID}
	if err := msg.Delete(); err != nil {
		context.Response(http.StatusBadRequest, 0, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, nil)
	}
}

// EditCrawlMsg
// @Summary	Edit crawl message
// @Param		id	path	int	true	"ID"	default(1)
// @Accept		json
// @Param		form	body	api.crawlMsgForm	true	"edit crawl message"
// @Produce	json
// @Success	200	{object}			http.Response
// @Failure	400	{object}			http.Response
// @Router		/api/v1/crawl-messages/{ID} 	[put]
func EditCrawlMsg(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		msgForm crawlMsgForm
		form    idForm
	)

	form.ID = com.StrTo(ctx.Param("id")).MustInt()
	httpCode, errCode := Valid(&form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	httpCode, errCode = BindAndValid(ctx, &msgForm)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	msg := &internal.CrawlMsg{
		DataType:        msgForm.DataType,
		DataId:          msgForm.DataId,
		AccountType:     msgForm.AccountType,
		TargetAccountId: msgForm.TargetAccountId,
	}

	if err := msg.Edit(); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, msg)
	}
}

// StartCrawlMsg
func StartCrawlMsg(ctx *gin.Context) {

}
