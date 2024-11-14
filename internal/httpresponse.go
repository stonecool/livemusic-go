package internal

import (
	"log"

	"github.com/gin-gonic/gin"
)

// 错误码常量
const (
	Success = iota
	Error
	ErrorInvalidParams
	ErrorNotExists
)

// Response 响应结构体
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Context 包装
type Context struct {
	*gin.Context
}

// 错误码消息映射
var returnMsg = map[int]string{
	Success:            "ok",
	Error:              "error",
	ErrorInvalidParams: "invalid params",
	ErrorNotExists:     "resource not exists",
}

// 获取错误消息
func getMsg(code int) string {
	if msg, ok := returnMsg[code]; ok {
		return msg
	}
	log.Printf("error code:%d not found", code)
	return returnMsg[Error]
}

// Response 响应方法
func (ctx *Context) Response(httpCode, errCode int, data interface{}) {
	ctx.JSON(httpCode, Response{
		Code: errCode,
		Msg:  getMsg(errCode),
		Data: data,
	})
}
