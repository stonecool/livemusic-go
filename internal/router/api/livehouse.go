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
