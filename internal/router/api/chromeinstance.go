package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal"
	http2 "github.com/stonecool/livemusic-go/internal/http"
	"net/http"
)

type chromeInstanceForm struct {
	Ip   string `json:"ip" valid:"Required;MaxSize(16)"`
	Port int    `json:"port" valid:"Required;Min(9222);Max(65535)"`
}

// CreateChromeInstance
// @Summary	Create a local chrome instance
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/create-instance [post]
func CreateChromeInstance(ctx *gin.Context) {
	context := http2.Context{Context: ctx}

	if ins, err := internal.CreateLocalChromeInstance(); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, ins)
	}
}

// BindChromeInstance
// @Summary	Bind a chrome instance
// @Accept		json
// @Param		form	body	api.chromeInstanceForm	true
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/bind-instance [post]
func BindChromeInstance(ctx *gin.Context) {
	var (
		context = http2.Context{Context: ctx}
		form    chromeInstanceForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != http2.Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if ins, err := internal.BindChromeInstance(form.Ip, form.Port); err != nil {
		context.Response(http.StatusBadRequest, http2.Error, nil)
	} else {
		context.Response(http.StatusCreated, http2.Success, ins)
	}
}

func GetChromeInstances() {

}

func DeleteChromeInstance() {

}

func EditChromeInstance() {

}
