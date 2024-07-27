package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal"
	http2 "github.com/stonecool/livemusic-go/internal/http"
	"github.com/unknwon/com"
	"net/http"
)

type crawlMsgForm struct {
	DataType    string `json:"data_type" valid:"Required;MaxSize(100)"`
	DataId      int    `json:"data_id" valid:"Required"`
	AccountType string `json:"account_type" valid:"Required;MaxSize(100)"`
	AccountId   string `json:"account_id" valid:"Required;MaxSize(100)"`
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
		DataType:    form.DataType,
		DataId:      form.DataId,
		AccountType: form.AccountType,
		AccountId:   form.AccountId,
	}

	if err := msg.Add(); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, msg)
	}
}

// GetCrawlMsg
// @Summary	Get a crawl message
// @Param		ID	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/crawl-messages/{ID} [delete]
func GetCrawlMsg(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    idForm
	)

	form.ID = com.StrTo(ctx.Param("ID")).MustInt()
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
// @Param		ID	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/crawl-messages/{ID} [delete]
func DeleteCrawlMsg(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    idForm
	)

	form.ID = com.StrTo(ctx.Param("ID")).MustInt()
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
// @Accept		json
// @Param		form	body	api.crawlMsgForm	true	"modify crawl message"
// @Param		ID	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}			http.Response
// @Failure	400	{object}			http.Response
// @Router		/api/v1/crawl-messages/{ID} 	[put]
func EditCrawlMsg(ctx *gin.Context) {
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

	idForm.ID = com.StrTo(ctx.Param("ID")).MustInt()
	httpCode, errCode = Valid(&idForm)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	msg := &internal.CrawlMsg{
		ID:          idForm.ID,
		DataType:    form.DataType,
		DataId:      form.DataId,
		AccountType: form.AccountType,
		AccountId:   form.AccountId,
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
