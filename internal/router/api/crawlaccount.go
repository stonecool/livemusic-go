package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal/account"
	"github.com/stonecool/livemusic-go/internal/chrome"
	"github.com/stonecool/livemusic-go/internal/client"
	"github.com/unknwon/com"
	"log"
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
		context = Context{Context: ctx}
		form    crawlAccountForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if a, err := account.CreateAccount(form.AccountType); err != nil {
		context.Response(http.StatusBadRequest, ErrorNotExists, nil)
	} else {
		context.Response(http.StatusCreated, Success, a)
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
		context = Context{Context: ctx}
		form    idForm
	)

	form.ID = com.StrTo(ctx.Param("id")).MustInt()
	httpCode, errCode := Valid(&form)
	if errCode != Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	account := &account.Account{ID: form.ID}
	if err := account.Get(); err != nil {
		context.Response(http.StatusBadRequest, 0, nil)
	} else {
		context.Response(http.StatusOK, 0, account)
	}
}

// GetCrawlAccounts
// @Summary	Get multiple accounts
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	500	{object}	http.Response
// @Router		/api/v1/crawl-accounts [get]
func GetCrawlAccounts(ctx *gin.Context) {
	var context = Context{Context: ctx}

	//account := &account.Account{}
	//if accounts, err := account.GetAll(); err != nil {
	context.Response(http.StatusBadRequest, 0, nil)
	//} else {
	//	context.Response(http.StatusBadRequest, 0, accounts)
	//}
}

// DeleteCrawlAccount
// @Summary	Delete crawl account
// @Param		id	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/crawl-accounts/{ID} [delete]
func DeleteCrawlAccount(ctx *gin.Context) {
	var (
		context = Context{Context: ctx}
		form    idForm
	)

	form.ID = com.StrTo(ctx.Param("id")).MustInt()
	httpCode, errCode := Valid(&form)
	if errCode != Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if err := account.DeleteAccount(form.ID); err != nil {
		context.Response(http.StatusBadRequest, 0, nil)
	} else {
		context.Response(http.StatusOK, 0, nil)
	}
}

// CrawlAccountWebSocket
// @Summary	Crawl Account websocket
// @Param		id	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	500	{object}	http.Response
// @Router		/api/v1/crawl-accounts/ws/{ID} [get]
func CrawlAccountWebSocket(ctx *gin.Context) {
	var (
		context = Context{Context: ctx}
		form    idForm
	)

	form.ID = com.StrTo(ctx.Param("id")).MustInt()
	httpCode, errCode := Valid(&form)
	if errCode != Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if err := client.HandleWebsocket(form.ID, ctx); err != nil {
		log.Printf("%v", err)
	}
}

// CrawlAccountLogin
// @Summary	Crawl Account websocket
// @Param		id	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	500	{object}	http.Response
// @Router		/api/v1/crawl-accounts/ws/{ID} [get]
func CrawlAccountLogin(ctx *gin.Context) {
	instanceID := 123
	accountType := "wx"

	chrome.GetPool().Login(instanceID, accountType)

}
