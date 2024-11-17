package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal"
	"github.com/unknwon/com"
	"net/http"
)

type livehouseForm struct {
	ID   int    `json:"ID" valid:"Required;Min(1)"`
	Name string `json:"name" valid:"Required;MaxSize(100)"`
}

// AddLivehouse
// @Summary	Add a livehouse
// @Accept		json
// @Param		form	body	api.livehouseForm	true	"created livehouse object"
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/livehouses [post]
func AddLivehouse(ctx *gin.Context) {
	var (
		context = Context{Context: ctx}
		form    livehouseForm
	)

	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	house := &internal.Livehouse{
		Name: form.Name,
	}

	if err := house.Add(); err != nil {
		context.Response(http.StatusBadRequest, Error, nil)
	} else {
		context.Response(http.StatusCreated, Success, house)
	}
}

// GetLivehouse
// @Summary	get a livehouse
// @Param		ID	path	int	true	"ID"	default(1)
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/livehouses/{ID} [get]
func GetLivehouse(ctx *gin.Context) {
	var (
		context = Context{Context: ctx}
		form    idForm
	)

	form.ID = com.StrTo(ctx.Param("ID")).MustInt()
	httpCode, errCode := Valid(&form)
	if errCode != Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	house := &internal.Livehouse{ID: form.ID}
	if err := house.Get(); err != nil {
		context.Response(http.StatusBadRequest, Error, nil)
	} else {
		context.Response(http.StatusCreated, Success, house)
	}
}

// GetLivehouses
// @Summary	Get all livehouse
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/livehouses [get]
func GetLivehouses(ctx *gin.Context) {
	var context = Context{Context: ctx}

	house := &internal.Livehouse{}
	if houses, err := house.GetAll(); err != nil {
		context.Response(http.StatusBadRequest, Error, nil)
	} else {
		context.Response(http.StatusCreated, Success, houses)
	}
}

// EditLivehouse
// @Summary	Edit a livehouse
// @Accept		json
// @Param		form	body	api.livehouseForm	true	"created livehouse object"
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/livehouses/{id} [put]
func EditLivehouse(ctx *gin.Context) {
	var (
		context = Context{Context: ctx}
		form    livehouseForm
	)

	form.ID = com.StrTo(ctx.Param("id")).MustInt()
	httpCode, errCode := BindAndValid(ctx, &form)
	if errCode != Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	house := internal.Livehouse{
		ID:   form.ID,
		Name: form.Name,
	}

	if err := house.Edit(); err != nil {
		context.Response(http.StatusBadRequest, Error, nil)
	} else {
		context.Response(http.StatusOK, Error, house)
	}
}

// DeleteLivehouse
// @Summary	delete a livehouse
// @Param		id path int true "ID" default(1)
// @Produce	json
// @Success	200	{object}	http.Response
// @Failure	400	{object}	http.Response
// @Router		/api/v1/livehouses/{id} [delete]
func DeleteLivehouse(ctx *gin.Context) {
	var (
		context = Context{Context: ctx}
		form    idForm
	)

	form.ID = com.StrTo(ctx.Param("id")).MustInt()
	httpCode, errCode := Valid(&form)
	if errCode != Success {
		context.Response(httpCode, errCode, nil)
		return
	}

	house := &internal.Livehouse{ID: form.ID}
	if err := house.Delete(); err != nil {
		context.Response(http.StatusBadRequest, Error, nil)
	} else {
		context.Response(http.StatusOK, Error, nil)
	}
}
