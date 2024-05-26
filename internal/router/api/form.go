package api

import (
	"github.com/beego/beego/v2/core/validation"
	"github.com/gin-gonic/gin"
	http2 "github.com/stonecool/livemusic-go/internal/http"
	"log"
	"net/http"
)

// BindAndValid binds and validates data
func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, http2.ErrorInvalidParams
	}

	return Valid(form)
}

func Valid(form interface{}) (int, int) {
	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, http2.Error
	}

	if !check {
		logErrors(valid.Errors)
		return http.StatusBadRequest, http2.ErrorInvalidParams
	}

	return http.StatusOK, http2.Success
}

// logErrors logs error logs
func logErrors(errors []*validation.Error) {
	for _, err := range errors {
		log.Println(err.Key, err.Message)
	}

	return
}
