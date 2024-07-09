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
	// 爬虫
	apiV1.POST("/crawls", api.AddCrawl)
	apiV1.GET("/crawls/:id", api.GetCrawl)
	apiV1.GET("/crawls", api.GetCrawls)
	apiV1.DELETE("/crawls/:id", api.DeleteCrawl)
	apiV1.GET("/crawls/ws/:id", api.CrawlWebSocket)

	// 爬虫消息生成者
	apiV1.POST("/msg-producers/:id", api.AddMsgProducer)
	apiV1.DELETE("/msg-producers/:id", api.DeleteCrawlMsg)
	apiV1.PUT("/msg-producers/start/:id", api.StartCrawlMsgProducer)
	apiV1.POST("/msg-producers/start/:id", api.StartCrawlMsgProducer)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return nil
}
