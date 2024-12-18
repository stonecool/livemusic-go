package api

import (
	"github.com/stonecool/livemusic-go/internal/task"
	"net/http"

	"github.com/gin-gonic/gin"
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
// @Param      form    body    crawlTaskForm    true    "created crawl task"
// @Produce    json
// @Success    200    {object}    Response
// @Failure    400    {object}    Response
// @Router     /api/v1/crawl-tasks     [post]
func AddCrawlTask(ctx *gin.Context) {
	var (
		context = Context{Context: ctx}
		form    crawlTaskForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if t, err := task.CreateTask(form.AccountType, form.DataType, form.DataId, form.CronSpec); err != nil {
		context.Response(http.StatusBadRequest, ErrorNotExists, nil)
	} else {
		context.Response(http.StatusCreated, Success, t)
	}
}
