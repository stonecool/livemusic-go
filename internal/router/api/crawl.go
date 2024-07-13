package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal"
	http2 "github.com/stonecool/livemusic-go/internal/http"
	"github.com/unknwon/com"
	"net/http"
	"reflect"
)

type addCAForm struct {
	AccountType string `json:"account_type" valid:"Required;MaxSize(255)"`
}

// AddCrawlAccount
//
//	@Summary	Add a crawl
//	@Accept		json
//	@Param		form	body	api.addCAForm	true	"created crawl account object"
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	400	{object}	http.Response
//	@Router		/api/v1/crawls [post]
func AddCrawlAccount(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    addCAForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if account, err := internal.AddCrawlAccount(form.AccountType); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, account)
	}
}

// GetCrawlAccount
//
//	@Summary	Get a single crawl
//	@Param		id	path	int	true	"ID"	default(1)
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	400	{object}	http.Response
//	@Router		/api/v1/crawls/{id} [get]
func GetCrawlAccount(ctx *gin.Context) {
	type getForm struct {
		ID int `valid:"Required;Min(1)"`
	}

	var (
		context = http2.Context{Context: ctx}
		form    getForm
		crawl   *internal.Account
	)

	form.ID = com.StrTo(ctx.Param("id")).MustInt()
	httpCode, errCode := Valid(&form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	c, err := internal.GetCrawlAccount(form.ID)
	if err != nil {
		context.Response(http.StatusBadRequest, 0, nil)
		return
	}

	if reflect.ValueOf(*crawl).IsZero() {
		context.Response(http.StatusOK, -1, nil)
		return
	}

	context.Response(http.StatusOK, 0, c)
}

// GetCrawlAccounts
//
//	@Summary	Get multiple accounts
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	500	{object}	http.Response
//	@Router		/api/v1/crawls [get]
func GetCrawlAccounts(ctx *gin.Context) {
}

// DeleteCrawlAccount
func DeleteCrawlAccount(ctx *gin.Context) {
}

func CrawlWebSocket(ctx *gin.Context) {
	//type Form struct {
	//	ID int `valid:"Required;Min(1)"`
	//}
	//
	//var (
	//	context = http2.Context{Context: ctx}
	//	form    Form
	//)
	//
	//form.ID = com.StrTo(ctx.Param("id")).MustInt()
	//httpCode, errCode := Valid(&form)
	//if errCode != http2.Success {
	//	context.Response(httpCode, errCode, nil)
	//	return
	//}
	//
	//account, err := internal.GetCrawlAccount(form.ID)
	//if err != nil {
	//	return
	//}
	//
	//client, err := NewClient(account, ctx)
	//if err != nil {
	//	return
	//}
	//
	//go client.Read()
	//go client.Write()
}
