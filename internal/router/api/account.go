package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/stonecool/livemusic-go/internal/crawl"
	http2 "github.com/stonecool/livemusic-go/internal/http"
	"github.com/unknwon/com"
	"go.uber.org/zap"
	"log"
	"net/http"
	"reflect"
	"time"
)

type addAccountForm struct {
	AccountType string `json:"account_type" valid:"Required;MaxSize(255)"`
}

// AddAccount
//
//	@Summary	Add a crawl account
//	@Accept		json
//	@Param		form	body	api.addAccountForm	true	"created account object"
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	400	{object}	http.Response
//	@Router		/api/v1/accounts [post]
func AddAccount(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    addAccountForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	accountType := form.AccountType
	_, ok := internal.AccountConfigMap[accountType]
	if !ok {
		context.Response(http.StatusBadRequest, http2.Error, nil)
		return
	}

	if account, err := crawl.AddAccount(form.AccountType); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, account)
	}
}

// GetAccount
//
//	@Summary	Get a single crawl account
//	@Param		id	path	int	true	"ID"	default(1)
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	400	{object}	http.Response
//	@Router		/api/v1/accounts/{id} [get]
func GetAccount(ctx *gin.Context) {
	type getForm struct {
		ID int `valid:"Required;Min(1)"`
	}

	var (
		context = http2.Context{Context: ctx}
		form    getForm
		account *crawl.Account
	)

	form.ID = com.StrTo(ctx.Param("id")).MustInt()
	httpCode, errCode := Valid(&form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	account, err := crawl.GetAccountByID(form.ID)
	if err != nil {
		context.Response(http.StatusBadRequest, 0, nil)
		return
	}

	if reflect.ValueOf(*account).IsZero() {
		context.Response(http.StatusOK, -1, nil)
		return
	}

	context.Response(http.StatusOK, 0, account)
}

// GetAccounts
//
//	@Summary	Get multiple accounts
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	500	{object}	http.Response
//	@Router		/api/v1/templates [get]
func GetAccounts(ctx *gin.Context) {

}

// DeleteAccount
func DeleteAccount(ctx *gin.Context) {

}

func CrawlWS(ctx *gin.Context) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			internal.Logger.Warn("defer ws connect error", zap.Error(err))
			return
		}
	}(conn)

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		err = conn.WriteMessage(mt, message)
		if err != nil {
			internal.Logger.Warn("ws write error", zap.Error(err))
			return
		}
		time.Sleep(time.Second)
	}
}
