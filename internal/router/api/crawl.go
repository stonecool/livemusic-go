package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal/crawl"
	http2 "github.com/stonecool/livemusic-go/internal/http"
	"github.com/unknwon/com"
	"net/http"
	"reflect"
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

	c, err := crawl.GetCrawlByID(form.ID)
	if err != nil {
		context.Response(http.StatusBadRequest, 0, nil)
		return
	}

	if reflect.ValueOf(*m).IsZero() {
		context.Response(http.StatusOK, -1, nil)
		return
	}

	context.Response(http.StatusOK, 0, c)
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

	c, err := crawl.GetCrawlByID(form.ID)
	if err != nil {
		return
	}

	client, err := crawl.NewClient(c, ctx)
	if err != nil {
		return
	}

	go client.Read()
	go client.Write()
}
