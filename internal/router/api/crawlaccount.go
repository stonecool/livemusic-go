package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal"
	crawl2 "github.com/stonecool/livemusic-go/internal/crawl"
	http2 "github.com/stonecool/livemusic-go/internal/http"
	"github.com/unknwon/com"
	"net/http"
)

type crawlAccountForm struct {
	AccountType string `json:"account_type" valid:"Required;MaxSize(255)"`
}

// AddCrawlAccount
// @Summary	Add crawl account
// @Accept		json
// @Param		form	body	api.crawlAccountForm	true	"created crawl account object"
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/crawl-accounts [post]
func AddCrawlAccount(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    crawlAccountForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	account := &internal.CrawlAccount{AccountType: form.AccountType}
	if err := account.Add(); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, account)
	}
}

// GetCrawlAccount
// @Summary	Get crawl account
// @Param		id	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/crawl-accounts/{id} [get]
func GetCrawlAccount(ctx *gin.Context) {
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

	account := &internal.CrawlAccount{ID: form.ID}
	if err := account.Get(); err != nil {
		context.Response(http.StatusBadRequest, 0, nil)
	} else {
		context.Response(http.StatusOK, 0, account)
	}
}

// GetCrawlAccounts
// @Summary	Get multiple accounts
// @Accept		json
// @Param		form	body	api.crawlAccountForm	true	"created crawl account object"
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	500	{object}	http.Response
// @Router		/api/v1/crawls [get]
func GetCrawlAccounts(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    crawlAccountForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	account := &internal.CrawlAccount{}
	if accounts, err := account.GetAll(); err != nil {
		context.Response(http.StatusBadRequest, 0, nil)
	} else {
		context.Response(http.StatusBadRequest, 0, accounts)
	}
}

// DeleteCrawlAccount
// @Summary	Delete crawl account
// @Param		id	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/crawl-accounts/{id} [delete]
func DeleteCrawlAccount(ctx *gin.Context) {
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

	account := &internal.CrawlAccount{ID: form.ID}
	if err := account.Delete(); err != nil {
		context.Response(http.StatusBadRequest, 0, nil)
	} else {
		context.Response(http.StatusOK, 0, nil)
	}
}

// CrawlAccountWebSocket
// @Summary	Get multiple accounts
// @Param		id	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	500	{object}	http.Response
// @Router		/api/v1/crawl-accounts/ws/{id} [get]
func CrawlAccountWebSocket(ctx *gin.Context) {
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

	crawl, err := crawl2.GetCrawl(form.ID)
	if err != nil {
		context.Response(http.StatusBadRequest, 0, nil)
		return
	}

	// TODO 同一个account,多个new的问题
	client, err := NewClient(crawl, ctx)
	if err != nil {
		return
	}

	go client.Read()
	go client.Write()

	context.Response(http.StatusOK, 0, nil)
}
