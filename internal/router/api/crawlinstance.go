package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal/crawl"
	http2 "github.com/stonecool/livemusic-go/internal/http"
	"github.com/unknwon/com"
	"net/http"
	"reflect"
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

	if _, err := crawl.GetCrawlAccountByID(form.AccountId); err != nil {
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
