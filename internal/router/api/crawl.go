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

type addCrawlForm struct {
	CrawlType string `json:"crawl_type" valid:"Required;MaxSize(255)"`
}

// AddCrawl
//
//	@Summary	Add a crawl
//	@Accept		json
//	@Param		form	body	api.addCrawlForm	true	"created crawl object"
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	400	{object}	http.Response
//	@Router		/api/v1/crawls [post]
func AddCrawl(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    addCrawlForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	crawlType := form.CrawlType
	_, ok := internal.CrawlAccountMap[crawlType]
	if !ok {
		context.Response(http.StatusBadRequest, http2.Error, nil)
		return
	}

	if account, err := crawl.AddCrawl(form.CrawlType); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, account)
	}
}

// GetCrawl
//
//	@Summary	Get a single crawl
//	@Param		id	path	int	true	"ID"	default(1)
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	400	{object}	http.Response
//	@Router		/api/v1/crawls/{id} [get]
func GetCrawl(ctx *gin.Context) {
	type getForm struct {
		ID int `valid:"Required;Min(1)"`
	}

	var (
		context = http2.Context{Context: ctx}
		form    getForm
		m       *crawl.Crawl
	)

	form.ID = com.StrTo(ctx.Param("id")).MustInt()
	httpCode, errCode := Valid(&form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	m, err := crawl.GetCrawlByID(form.ID)
	if err != nil {
		context.Response(http.StatusBadRequest, 0, nil)
		return
	}

	if reflect.ValueOf(*m).IsZero() {
		context.Response(http.StatusOK, -1, nil)
		return
	}

	context.Response(http.StatusOK, 0, m)
}

// GetCrawls
//
//	@Summary	Get multiple accounts
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	500	{object}	http.Response
//	@Router		/api/v1/crawls [get]
func GetCrawls(ctx *gin.Context) {
}

// DeleteCrawl
func DeleteCrawl(ctx *gin.Context) {
}

type WebsocketClient struct {
	crawl *crawl.ICrawl
	conn  websocket.Conn
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
