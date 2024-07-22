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
	apiV1.GET("/crawl-accounts/ws/:id", api.CrawlAccountWebSocket)

	// 爬虫消息
	apiV1.POST("/crawl-messages/:id", api.AddCrawlMsg)
	apiV1.GET("/crawl-messages/:id", api.GetCrawlMsg)
	apiV1.GET("/crawl-messages", api.GetCrawlMsgs)
	apiV1.DELETE("/crawl-messages/:id", api.DeleteCrawlMsg)
	apiV1.PUT("/crawl-messages/start/:id", api.ModifyCrawlMsg)
	apiV1.POST("/crawl-messages/start/:id", api.StartCrawlMsgProducer)

	// livehouse
	apiV1.POST("/livehouses", api.AddLivehouse)
	apiV1.DELETE("/livehouses/:id", api.DeleteLivehouse)
	apiV1.PUT("/livehouses/:id", api.EditLivehouse)
	apiV1.GET("/livehouses/:id", api.GetLivehouse)
	apiV1.GET("/livehouses", api.GetLivehouses)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return nil
}
