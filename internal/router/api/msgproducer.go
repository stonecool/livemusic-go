package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal"
	http2 "github.com/stonecool/livemusic-go/internal/http"
	"net/http"
)

type addMsgProducerForm struct {
	DataType  string `json:"data_type" valid:"Required;MaxSize(100)"`
	DataId    int    `json:"data_id" valid:"Required"`
	CrawlType string `json:"crawl_type" valid:"Required;MaxSize(100)"`
	AccountId string `json:"account_id" valid:"Required;MaxSize(100)"`
}

// AddMsgProducer
//
//	@Summary	Adds a crawl message producer
//	@Accept		json
//	@Param		form	body	api.addMsgProducerForm	true	"created crawl message producer"
//	@Produce	json
//	@Success	200	{object}			http.Response
//	@Failure	400	{object}			http.Response
//	@Router		/api/v1/msg-producers 	[post]
func AddMsgProducer(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    addMsgProducerForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if msgProducer, err := internal.AddMsgProducer(form.DataType, form.DataId, form.CrawlType, form.AccountId); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, msgProducer)
	}
}

// DeleteCrawlMsg
//
//	@Summary	Delete a crawl message producer
//	@Param		id	path	int	true	"ID"	default(1)
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	400	{object}	http.Response
//	@Router		/api/v1/msg-producers/{id} [get]
func DeleteCrawlMsg(ctx *gin.Context) {
	//type deleteForm struct {
	//	ID int `valid:"Required;Min(1)"`
	//}
	//
	//var (
	//	context = http2.Context{Context: ctx}
	//	form    deleteForm
	//)
	//
	//form.ID = com.StrTo(ctx.Param("id")).MustInt()
	//httpCode, errCode := Valid(&form)
	//if errCode != http2.Success {
	//	context.Response(httpCode, errCode, nil)
	//	return
	//}
	//
	//c, err := crawl.GetCrawlAccount(form.ID)
	//if err != nil {
	//	context.Response(http.StatusBadRequest, 0, nil)
	//	return
	//}
	//
	//if reflect.ValueOf(*m).IsZero() {
	//	context.Response(http.StatusOK, -1, nil)
	//	return
	//}
	//
	//context.Response(http.StatusOK, 0, c)
}

// ModifyCrawlMsg
func ModifyCrawlMsg(ctx *gin.Context) {

}

// StartCrawlMsgProducer
func StartCrawlMsgProducer(ctx *gin.Context) {

}
