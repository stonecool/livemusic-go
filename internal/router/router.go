package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/stonecool/livemusic-go/docs/swagger"
	"github.com/stonecool/livemusic-go/internal/router/api"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(engine *gin.Engine) error {
	router := engine.Group("/")

	apiV1 := router.Group("/api/v1")
	// 爬虫账号
	apiV1.POST("/accounts", api.AddAccount)
	apiV1.GET("/accounts/:id", api.GetAccount)
	apiV1.GET("/accounts", api.GetAccounts)
	apiV1.DELETE("/accounts/:id", api.DeleteAccount)
	apiV1.GET("/crawls/ws/:id", api.CrawlWS)

	// 爬虫实例
	apiV1.POST("/crawlInstances", api.AddCrawlInstance)
	apiV1.GET("/crawlInstances/:id", api.GetCrawlInstance)
	apiV1.GET("/crawlInstances", api.GetCrawlInstances)
	apiV1.PUT("/crawlInstances/:id", api.EditCrawlInstance)
	apiV1.DELETE("/crawlInstances/:id", api.DeleteCrawlInstance)

	// 爬虫
	apiV1.GET("/crawl/:id", api.EditCrawlInstance)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return nil
}
