package http

import "log"

const (
	Success = iota
	Error
	ErrorInvalidParams
	ErrorNotExists
)

var returnMsg = map[int]string{
	Success:            "ok",
	Error:              "error",
	ErrorInvalidParams: "invalid params",
	ErrorNotExists:     "resource not exists",
}

func getMsg(code int) string {
	if msg, ok := returnMsg[code]; ok {
		return msg
	} else {
		log.Printf("error code:%d not found", code)
		return returnMsg[Error]
	}
}
