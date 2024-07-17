package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal"
	http2 "github.com/stonecool/livemusic-go/internal/http"
	"net/http"
)

type addLivehouseForm struct {
	Name string `json:"name" valid:"Required;MaxSize(100)"`
}

// AddLivehouse
//
//	@Summary	Add a livehouse
//	@Accept		json
//	@Param		form	body	api.addLivehouseForm	true	"created livehouse object"
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	400	{object}	http.Response
//	@Router		/api/v1/livehouses [post]
func AddLivehouse(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    addLivehouseForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if account, err := internal.AddCrawlAccount(""); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, account)
	}
}

// DeleteLivehouse
//
//	@Summary	delete a livehouse
//	@Param		id path int true "ID" default(1)
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	400	{object}	http.Response
//	@Router		/api/v1/livehouses/{id} [delete]
func DeleteLivehouse(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    addLivehouseForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if account, err := internal.AddCrawlAccount(""); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, account)
	}
}

// ModifyLivehouse
//
//	@Summary	modify a livehouse
//	@Accept		json
//	@Param		form	body	api.addLivehouseForm	true	"created livehouse object"
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	400	{object}	http.Response
//	@Router		/api/v1/livehouses/{id} [post]
func ModifyLivehouse(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    addLivehouseForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if account, err := internal.AddCrawlAccount(""); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, account)
	}
}

// GetLivehouse
//
//	@Summary	get a livehouse
//	@Param		id	path	int	true	"ID"	default(1)
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	400	{object}	http.Response
//	@Router		/api/v1/livehouses/{id} [get]
func GetLivehouse(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    addLivehouseForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if account, err := internal.AddCrawlAccount(""); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, account)
	}
}

// GetLivehouses
//
//	@Summary	Add a livehouse
//	@Accept		json
//	@Param		form	body	api.addLivehouseForm	true	"created livehouse object"
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	400	{object}	http.Response
//	@Router		/api/v1/livehouses [post]
func GetLivehouses(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    addLivehouseForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if account, err := internal.AddCrawlAccount(""); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, account)
	}
}
