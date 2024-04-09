package api

import (
	"github.com/gin-gonic/gin"
)

// InitCrawl
//
//	@Summary	Get a single crawl instance
//	@Param		id	path	int	true	"ID"	default(1)
//	@Produce	json
//	@Success	200	{object}	http.Response
//	@Failure	400	{object}	http.Response
//	@Router		/api/v1/crawl/{id} [get]
func InitCrawl(ctx *gin.Context) {
	//id := 1
	//
	//if _, err := bak.GetCrawlByID(id); err == nil {
	//}
	//
	//bak.InitCrawl(id)
}
