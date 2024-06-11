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

type addInstanceForm struct {
	Name        string                 `json:"name" valid:"Required;MaxSize(255)"`
	AccountId   int                    `json:"account_id" valid:"Required;Min(1)"`
	Headers     map[string]interface{} `json:"headers"`
	QueryParams string                 `json:"query_params" valid:"MaxSize(255)"`
	FormData    string                 `json:"form_data" valid:"MaxSize(255)"`
}

// AddCrawlInstance
//
//	@Summary	Add crawl instance
//	@Accept		json
//	@Param		instance	body	api.addInstanceForm	true	"created instance object"
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	400	{object}	http.Response
//	@Router		/api/v1/crawlInstances [post]
func AddCrawlInstance(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    addInstanceForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if _, err := crawl.GetAccountByID(form.AccountId); err != nil {
		context.Response(http.StatusBadRequest, http2.ErrorNotExists, nil)
		return
	}

	instance := crawl.Instance{
		Name:      form.Name,
		AccountId: form.AccountId,
	}

	if err := instance.Add(); err != nil {
		context.Response(http.StatusBadRequest, -1, nil)
	} else {
		context.Response(http.StatusCreated, 0, instance)
	}
}

// GetCrawlInstance
//
//	@Summary	Get a single crawl instance
//	@Param		id	path	int	true	"ID"	default(1)
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	400	{object}	http.Response
//	@Router		/api/v1/crawlInstances/{id} [get]
func GetCrawlInstance(ctx *gin.Context) {
	type Form struct {
		ID int `valid:"Required;Min(1)"`
	}

	var (
		context  = http2.Context{Context: ctx}
		form     Form
		instance *crawl.Instance
	)

	form.ID = com.StrTo(ctx.Param("id")).MustInt()
	httpCode, errCode := Valid(&form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	instance, err := crawl.GetCrawlInstanceByID(form.ID)
	if err != nil {
		context.Response(http.StatusBadRequest, 0, nil)
		return
	}

	if reflect.ValueOf(*instance).IsZero() {
		context.Response(http.StatusOK, -1, nil)
		return
	}

	context.Response(http.StatusOK, 0, instance)
}

// GetCrawlInstances
func GetCrawlInstances(ctx *gin.Context) {

}

// EditCrawlInstance
func EditCrawlInstance(ctx *gin.Context) {

}

// DeleteCrawlInstance
func DeleteCrawlInstance(ctx *gin.Context) {

}

func CrawlWS(ctx *gin.Context) {
	type Form struct {
		ID int `valid:"Required;Min(1)"`
	}

	var (
		context = http2.Context{Context: ctx}
		form    Form
	)

	form.ID = com.StrTo(ctx.Param("id")).MustInt()
	httpCode, errCode := Valid(&form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

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
