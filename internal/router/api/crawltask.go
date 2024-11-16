package api

import (
	"github.com/stonecool/livemusic-go/internal/task"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal"
)

type crawlTaskForm struct {
	DataType        string `json:"data_type" valid:"Required;MaxSize(100)"`
	DataId          int    `json:"data_id" valid:"Required"`
	AccountType     string `json:"account_type" valid:"Required;MaxSize(100)"`
	TargetAccountId string `json:"target_account_id" valid:"Required;MaxSize(100)"`
	CronSpec        string `json:"cron_spec" valid:"Required"`
}

// @Summary    Adds crawl task
// @Accept     json
// @Param      form    body    api.crawlTaskForm    true    "created crawl task"
// @Produce    json
// @Success    200    {object}    http.Response
// @Failure    400    {object}    http.Response
// @Router     /api/v1/crawl-tasks     [post]
func AddCrawlTask(ctx *gin.Context) {
	var (
		context = internal.Context{Context: ctx}
		form    crawlTaskForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != internal.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	task := task.Task{
		DataType:        form.DataType,
		DataId:          form.DataId,
		AccountType:     form.AccountType,
		TargetAccountId: form.TargetAccountId,
		CronSpec:        form.CronSpec,
	}

	if err := task.Add(); err != nil {
		context.Response(http.StatusBadRequest, internal.ErrorNotExists, nil)
	} else {
		context.Response(http.StatusCreated, internal.Success, task)
	}
}
