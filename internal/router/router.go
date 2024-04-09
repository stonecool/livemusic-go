package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/stonecool/1701livehouse-server/docs/swagger"
	"github.com/stonecool/1701livehouse-server/internal/router/api"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(engine *gin.Engine) error {
	router := engine.Group("/")

	apiV1 := router.Group("/api/v1")
	// 爬虫实例模板
	apiV1.POST("/crawlAccounts", api.AddCrawlAccount)
	apiV1.GET("/crawlAccounts/:id", api.GetCrawlAccount)
	apiV1.GET("/crawlAccounts", api.GetCrawlAccounts)
	apiV1.PUT("/crawlAccounts/:id", api.EditCrawlAccount)
	apiV1.DELETE("/crawlAccounts/:id", api.DeleteCrawlAccount)

	// 爬虫实例
	apiV1.POST("/crawlInstances", api.AddCrawlInstance)
	apiV1.GET("/crawlInstances/:id", api.GetCrawlInstance)
	apiV1.GET("/crawlInstances", api.GetCrawlInstances)
	apiV1.PUT("/crawlInstances/:id", api.EditCrawlInstance)
	apiV1.DELETE("/crawlInstances/:id", api.DeleteCrawlInstance)

	// 爬虫
	apiV1.GET("/crawl/:id", api.InitCrawl)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return nil
}
