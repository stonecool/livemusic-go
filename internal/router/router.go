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
	apiV1.POST("/crawl-accounts", api.AddCrawlAccount)
	apiV1.GET("/crawl-accounts/:id", api.GetCrawlAccount)
	apiV1.GET("/crawl-accounts", api.GetCrawlAccounts)
	apiV1.DELETE("/crawl-accounts/:id", api.DeleteCrawlAccount)
	apiV1.GET("/crawl-accounts/ws/:id", api.CrawlWebSocket)

	// 爬虫消息
	apiV1.POST("/crawl-messages/:id", api.AddCrawlMsg)
	apiV1.DELETE("/crawl-messages/:id", api.DeleteCrawlMsg)
	apiV1.PUT("/crawl-messages/start/:id", api.StartCrawlMsgProducer)
	apiV1.POST("/crawl-messages/start/:id", api.StartCrawlMsgProducer)

	// livehouse
	
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return nil
}
