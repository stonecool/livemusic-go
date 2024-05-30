package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal/crawl"
	http2 "github.com/stonecool/livemusic-go/internal/http"
	"github.com/unknwon/com"
	"net/http"
)

type addAccountForm struct {
	AccountType string `json:"account_type" valid:"Required;MaxSize(255)"`
}

// AddCrawlAccount
//
//	@Summary	Add crawl template
//	@Accept		json
//	@Param		template	body	api.addAccountForm	true	"created template object"
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	400	{object}	http.Response
//	@Router		/api/v1/crawlAccounts [post]
func AddCrawlAccount(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    addAccountForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	account := crawl.Account{
		AccountType: form.AccountType,
	}

	if err := account.Add(); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, account)
	}
}

// GetCrawlAccount
//
//	@Summary	Get a single crawl account
//	@Param		id	path	int	true	"ID"	default(1)
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	400	{object}	http.Response
//	@Router		/api/v1/crawlAccounts/{id} [get]
func GetCrawlAccount(ctx *gin.Context) {
	type Form struct {
		ID int `valid:"Required;Min(1)"`
	}

	var (
		context  = http2.Context{Context: ctx}
		form     Form
		template *crawl.Account
	)

	form.ID = com.StrTo(ctx.Param("id")).MustInt()
	httpCode, errCode := Valid(&form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	//template, err := crawl.GetCrawlAccountByID(form.ID)
	//if err != nil {
	//	context.Response(http.StatusBadRequest, 0, nil)
	//	return
	//}
	//
	//if reflect.ValueOf(*template).IsZero() {
	//	context.Response(http.StatusOK, -1, nil)
	//	return
	//}
	//
	context.Response(http.StatusOK, 0, template)
}

// GetCrawlAccounts
//
//	@Summary	Get multiple accounts
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	500	{object}	http.Response
//	@Router		/api/v1/templates [get]
func GetCrawlAccounts(ctx *gin.Context) {

}

// EditCrawlAccount
func EditCrawlAccount(context *gin.Context) {

}

// DeleteCrawlAccount
func DeleteCrawlAccount(context *gin.Context) {

}
