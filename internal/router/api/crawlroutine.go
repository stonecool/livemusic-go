package api

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
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
// @Param		form	body	crawlRoutineForm	true	"created crawl message"
// @Produce	json
// @Success	200	{object}			Response
// @Failure	400	{object}			Response
// @Router		/api/v1/crawl-routines 	[post]
func AddCrawlRoutine(ctx *gin.Context) {
	var (
		context = Context{Context: ctx}
		form    crawlRoutineForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	//routine := internal.CrawlRoutine{
	//	DataType:        form.DataType,
	//	DataId:          form.DataId,
	//	AccountType:     form.AccountType,
	//	TargetAccountId: form.TargetAccountId,
	//}
	//
	//if err := routine.Add(); err != nil {
	//	context.Response(http.StatusBadRequest, Error, nil)
	//} else {
	//	context.Response(http.StatusCreated, Success, routine)
	//}
}

// GetCrawlRoutine
// @Summary	Get a crawl routine
// @Param		id	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}	Response
// @Failure	400	{object}	Response
// @Router		/api/v1/crawl-routines/{ID} [get]
func GetCrawlRoutine(ctx *gin.Context) {
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

	//routine := internal.CrawlRoutine{ID: form.ID}
	//if err := routine.Get(); err != nil {
	//	context.Response(http.StatusBadRequest, 0, nil)
	//} else {
	//	context.Response(http.StatusCreated, Success, routine)
	//}
}

// GetCrawlRoutines
// @Summary	Get all crawl routines
// @Produce	json
// @Success	200	{object}	Response
// @Failure	400	{object}	Response
// @Router		/api/v1/crawl-routines [get]
func GetCrawlRoutines(ctx *gin.Context) {
	//var context = Context{Context: ctx}

	//message := internal.CrawlRoutine{}
	//if msgs, err := message.GetAll(); err != nil {
	//	context.Response(http.StatusBadRequest, 0, nil)
	//} else {
	//	context.Response(http.StatusCreated, Success, msgs)
	//}
}

// DeleteCrawlRoutine
// @Summary	Delete crawl routine
// @Param		id	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}	Response
// @Failure	400	{object}	Response
// @Router		/api/v1/crawl-routines/{ID} [delete]
func DeleteCrawlRoutine(ctx *gin.Context) {
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

	//routine := &internal.CrawlRoutine{ID: form.ID}
	//if err := routine.Delete(); err != nil {
	//	context.Response(http.StatusBadRequest, 0, nil)
	//} else {
	//	context.Response(http.StatusCreated, Success, nil)
	//}
}

// EditCrawlRoutine
// @Summary	Edit crawl routine
// @Param		id	path	int	true	"ID"	default(1)
// @Accept		json
// @Param		form	body	crawlRoutineForm	true	"edit crawl message"
// @Produce	json
// @Success	200	{object}			Response
// @Failure	400	{object}			Response
// @Router		/api/v1/crawl-routines/{ID} 	[put]
func EditCrawlRoutine(ctx *gin.Context) {
	var (
		context = Context{Context: ctx}
		msgForm crawlRoutineForm
		form    idForm
	)

	form.ID = com.StrTo(ctx.Param("id")).MustInt()
	httpCode, errCode := Valid(&form)
	if errCode != Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	httpCode, errCode = BindAndValid(ctx, &msgForm)
	if errCode != Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	//message := &internal.CrawlRoutine{
	//	DataType:        msgForm.DataType,
	//	DataId:          msgForm.DataId,
	//	AccountType:     msgForm.AccountType,
	//	TargetAccountId: msgForm.TargetAccountId,
	//}
	//
	//if err := message.Edit(); err != nil {
	//	context.Response(http.StatusBadRequest, Error, nil)
	//} else {
	//	context.Response(http.StatusCreated, Success, message)
	//}
}

// StartCrawlRoutine
// @Summary	Start a crawl message
// @Param		id	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}	Response
// @Failure	400	{object}	Response
// @Router		/api/v1/crawl-routines/start/{ID} [post]
func StartCrawlRoutine(ctx *gin.Context) {

}
