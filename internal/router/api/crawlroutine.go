package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/unknwon/com"
	"net/http"
)

type crawlRoutineForm struct {
	DataType        string `json:"data_type" valid:"Required;MaxSize(100)"`
	DataId          int    `json:"data_id" valid:"Required"`
	AccountType     string `json:"account_type" valid:"Required;MaxSize(100)"`
	TargetAccountId string `json:"target_account_id" valid:"Required;MaxSize(100)"`
}

// AddCrawlRoutine
// @Summary	Adds crawl routine
// @Accept		json
// @Param		form	body	api.crawlRoutineForm	true	"created crawl message"
// @Produce	json
// @Success	200	{object}			http.Response
// @Failure	400	{object}			http.Response
// @Router		/api/v1/crawl-routines 	[post]
func AddCrawlRoutine(ctx *gin.Context) {
	var (
		context = internal.Context{Context: ctx}
		form    crawlRoutineForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != internal.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	routine := internal.CrawlRoutine{
		DataType:        form.DataType,
		DataId:          form.DataId,
		AccountType:     form.AccountType,
		TargetAccountId: form.TargetAccountId,
	}

	if err := routine.Add(); err != nil {
		context.Response(http.StatusBadRequest, internal.Error, nil)
	} else {
		context.Response(http.StatusCreated, internal.Success, routine)
	}
}

// GetCrawlRoutine
// @Summary	Get a crawl routine
// @Param		id	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/crawl-routines/{ID} [get]
func GetCrawlRoutine(ctx *gin.Context) {
	var (
		context = internal.Context{Context: ctx}
		form    idForm
	)

	form.ID = com.StrTo(ctx.Param("id")).MustInt()
	httpCode, errCode := Valid(&form)
	if errCode != internal.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	routine := internal.CrawlRoutine{ID: form.ID}
	if err := routine.Get(); err != nil {
		context.Response(http.StatusBadRequest, 0, nil)
	} else {
		context.Response(http.StatusCreated, internal.Success, routine)
	}
}

// GetCrawlRoutines
// @Summary	Get all crawl routines
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/crawl-routines [get]
func GetCrawlRoutines(ctx *gin.Context) {
	var context = internal.Context{Context: ctx}

	msg := internal.CrawlRoutine{}
	if msgs, err := msg.GetAll(); err != nil {
		context.Response(http.StatusBadRequest, 0, nil)
	} else {
		context.Response(http.StatusCreated, internal.Success, msgs)
	}
}

// DeleteCrawlRoutine
// @Summary	Delete crawl routine
// @Param		id	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/crawl-routines/{ID} [delete]
func DeleteCrawlRoutine(ctx *gin.Context) {
	var (
		context = internal.Context{Context: ctx}
		form    idForm
	)

	form.ID = com.StrTo(ctx.Param("id")).MustInt()
	httpCode, errCode := Valid(&form)
	if errCode != internal.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	routine := &internal.CrawlRoutine{ID: form.ID}
	if err := routine.Delete(); err != nil {
		context.Response(http.StatusBadRequest, 0, nil)
	} else {
		context.Response(http.StatusCreated, internal.Success, nil)
	}
}

// EditCrawlRoutine
// @Summary	Edit crawl routine
// @Param		id	path	int	true	"ID"	default(1)
// @Accept		json
// @Param		form	body	api.crawlRoutineForm	true	"edit crawl message"
// @Produce	json
// @Success	200	{object}			http.Response
// @Failure	400	{object}			http.Response
// @Router		/api/v1/crawl-routines/{ID} 	[put]
func EditCrawlRoutine(ctx *gin.Context) {
	var (
		context = internal.Context{Context: ctx}
		msgForm crawlRoutineForm
		form    idForm
	)

	form.ID = com.StrTo(ctx.Param("id")).MustInt()
	httpCode, errCode := Valid(&form)
	if errCode != internal.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	httpCode, errCode = BindAndValid(ctx, &msgForm)
	if errCode != internal.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	msg := &internal.CrawlRoutine{
		DataType:        msgForm.DataType,
		DataId:          msgForm.DataId,
		AccountType:     msgForm.AccountType,
		TargetAccountId: msgForm.TargetAccountId,
	}

	if err := msg.Edit(); err != nil {
		context.Response(http.StatusBadRequest, internal.Error, nil)
	} else {
		context.Response(http.StatusCreated, internal.Success, msg)
	}
}

// StartCrawlRoutine
// @Summary	Start a crawl msg
// @Param		id	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/crawl-routines/start/{ID} [pos]
func StartCrawlRoutine(ctx *gin.Context) {

}
