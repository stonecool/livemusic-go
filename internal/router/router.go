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
	// chrome实例
	apiV1.POST("/chromes", api.CreateChrome)
	apiV1.POST("/chromes/bind", api.BindChrome)
	apiV1.GET("/chromes", api.ListChromes)
	apiV1.GET("/chromes/:id", api.GetChrome)
	apiV1.DELETE("/chromes/:id", api.DeleteChrome)

	// 爬虫账号
	apiV1.POST("/crawl-accounts", api.AddCrawlAccount)
	apiV1.GET("/crawl-accounts/:id", api.GetCrawlAccount)
	apiV1.GET("/crawl-accounts", api.GetCrawlAccounts)
	apiV1.DELETE("/crawl-accounts/:id", api.DeleteCrawlAccount)
	apiV1.GET("/crawl-accounts/ws/:id", api.CrawlAccountWebSocket)

	// 爬虫例程
	apiV1.POST("/crawl-routines", api.AddCrawlRoutine)
	apiV1.GET("/crawl-routines/:id", api.GetCrawlRoutine)
	apiV1.GET("/crawl-routines", api.GetCrawlRoutines)
	apiV1.DELETE("/crawl-routines/:id", api.DeleteCrawlRoutine)
	apiV1.PUT("/crawl-routines/start/:id", api.EditCrawlRoutine)
	apiV1.POST("/crawl-routines/start/:id", api.StartCrawlRoutine)

	// livehouse
	apiV1.POST("/livehouses", api.AddLivehouse)
	apiV1.DELETE("/livehouses/:id", api.DeleteLivehouse)
	apiV1.PUT("/livehouses/:id", api.EditLivehouse)
	apiV1.GET("/livehouses/:id", api.GetLivehouse)
	apiV1.GET("/livehouses", api.GetLivehouses)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return nil
}
