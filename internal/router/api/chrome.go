package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal/chrome"
	"net/http"
)

type chromeForm struct {
	Ip   string `json:"ip" valid:"Required;MaxSize(16)"`
	Port int    `json:"port" valid:"Required;Min(9222);Max(65535)"`
}

// CreateChrome
// @Summary	Create a local chrome instance
// @Produce	json
// @Param		form	body	chromeForm	true "form"
// @Success	200	{object}	Response
// @Failure	400	{object}	Response
// @Router		/api/v1/create-instance [post]
func CreateChrome(ctx *gin.Context) {
	context := Context{Context: ctx}

	if ins, err := chrome.Create(); err != nil {
		context.Response(http.StatusBadRequest, Error, nil)
	} else {
		context.Response(http.StatusCreated, Success, ins)
	}
}

// BindChrome
// @Summary	Bind a chrome instance
// @Accept		json
// @Param		form	body	chromeForm	true "form"
// @Produce	json
// @Success	200	{object}	Response
// @Failure	400	{object}	Response
// @Router		/api/v1/bind-instance [post]
func BindChrome(ctx *gin.Context) {
	var (
		context = Context{Context: ctx}
		form    chromeForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if ins, err := chrome.Bind(form.Ip, form.Port); err != nil {
		context.Response(http.StatusBadRequest, Error, nil)
	} else {
		context.Response(http.StatusCreated, Success, ins)
	}
}

// GetChrome
// @Summary	Get multiple chrome instances
// @Produce	json
// @Success	200	{object}	Response
// @Failure	500	{object}	Response
// @Router		/api/v1/instances [get]
func GetChrome(ctx *gin.Context) {
	var context = Context{Context: ctx}

	if chromes, err := chrome.GetAll(); err != nil {
		context.Response(http.StatusBadRequest, 0, nil)
	} else {
		context.Response(http.StatusBadRequest, 0, chromes)
	}
}

func DeleteChrome(ctx *gin.Context) {

}
