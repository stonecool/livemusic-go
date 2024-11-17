package api

import (
	"github.com/beego/beego/v2/core/validation"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type idForm struct {
	ID int `valid:"Required;Min(1)"`
}

// BindAndValid binds and validates data
func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, ErrorInvalidParams
	}

	return Valid(form)
}

func Valid(form interface{}) (int, int) {
	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, Error
	}

	if !check {
		logErrors(valid.Errors)
		return http.StatusBadRequest, ErrorInvalidParams
	}

	return http.StatusOK, Success
}

// logErrors logs error logs
func logErrors(errors []*validation.Error) {
	for _, err := range errors {
		log.Println(err.Key, err.Message)
	}

	return
}
