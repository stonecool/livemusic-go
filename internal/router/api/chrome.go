package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal/chrome"
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
	context := Context{Context: ctx}

	if ins, err := chrome.CreateLocalChromeInstance(); err != nil {
		context.Response(http.StatusBadRequest, Error, nil)
	} else {
		context.Response(http.StatusCreated, Success, ins)
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
		context = Context{Context: ctx}
		form    chromeInstanceForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if ins, err := chrome.BindChromeInstance(form.Ip, form.Port); err != nil {
		context.Response(http.StatusBadRequest, Error, nil)
	} else {
		context.Response(http.StatusCreated, Success, ins)
	}
}

// GetChromeInstances
// @Summary	Get multiple chrome instances
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	500	{object}	http.Response
// @Router		/api/v1/instances [get]
func GetChromeInstances(ctx *gin.Context) {
	var context = Context{Context: ctx}

	if chromes, err := chrome.GetAllChrome(); err != nil {
		context.Response(http.StatusBadRequest, 0, nil)
	} else {
		context.Response(http.StatusBadRequest, 0, chromes)
	}
}

func DeleteChromeInstance(ctx *gin.Context) {

}