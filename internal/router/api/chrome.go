package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal/chrome"
)

type instanceForm struct {
	IP   string `json:"ip" valid:"Required;MaxSize(16)"`
	Port int    `json:"port" valid:"Required;Min(9222);Max(65535)"`
}

// CreateChrome
// @Summary Create a local new chrome
// @Produce json
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Router  /api/v1/chromes [post]
func CreateChrome(ctx *gin.Context) {
	context := Context{Context: ctx}

	if ins, err := chrome.Create(); err != nil {
		context.Response(http.StatusBadRequest, Error, nil)
	} else {
		context.Response(http.StatusCreated, Success, ins)
	}
}

// BindChrome
// @Summary Bind an existing chrome
// @Accept  json
// @Param   form    body    instanceForm    true "Instance configuration"
// @Produce json
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Router  /api/v1/chromes/bind [post]
func BindChrome(ctx *gin.Context) {
	var (
		context = Context{Context: ctx}
		form    instanceForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	if ins, err := chrome.Bind(form.IP, form.Port); err != nil {
		context.Response(http.StatusBadRequest, Error, nil)
	} else {
		context.Response(http.StatusCreated, Success, ins)
	}
}

// GetChrome
// @Summary Get a browser instance by ID
// @Produce json
// @Param   id     path    int    true "Instance ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router  /api/v1/chromes/{id} [get]
func GetChrome(ctx *gin.Context) {
	var context = Context{Context: ctx}

	//id := ctx.Param("id")
	//if id == "" {
	//	context.Response(http.StatusBadRequest, Error, "invalid instance id")
	//	return
	//}
	//
	//instance, err := chrome.GetByID(id)
	//if err != nil {
	//	if err == chrome.ErrInstanceNotFound {
	//		context.Response(http.StatusNotFound, Error, "instance not found")
	//	} else {
	//		context.Response(http.StatusInternalServerError, Error, "failed to get instance")
	//	}
	//	return
	//}

	context.Response(http.StatusOK, Success, nil)
}

// ListChromes
// @Summary List all browser instances
// @Produce json
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router  /api/v1/chromes [get]
func ListChromes(ctx *gin.Context) {
	var context = Context{Context: ctx}

	if instances, err := chrome.GetAll(); err != nil {
		context.Response(http.StatusInternalServerError, Error, nil)
	} else {
		context.Response(http.StatusOK, Success, instances)
	}
}

// DeleteChrome
// @Summary Delete a browser instance
// @Produce json
// @Param   id     path    int    true "Instance ID"
// @Success 204 {object} Response
// @Failure 404 {object} Response
// @Router  /api/v1/chromes/{id} [delete]
func DeleteChrome(ctx *gin.Context) {
	var context = Context{Context: ctx}
	// TODO: 实现删除实例的逻辑
	context.Response(http.StatusNoContent, Success, nil)
}
