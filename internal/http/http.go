package http

import (
	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (ctx *Context) Response(httpCode, errCode int, data interface{}) {
	ctx.JSON(httpCode, Response{
		Code: errCode,
		Msg:  getMsg(errCode),
		Data: data,
	})
}
