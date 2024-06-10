package internal

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

var Logger *zap.Logger
var SugarLogger *zap.SugaredLogger

func init() {
	var err error
	if gin.Mode() == gin.DebugMode {
		Logger, err = zap.NewDevelopment()
	} else {
		Logger, err = zap.NewProduction()
	}

	if err != nil {
		log.Fatalf("init zap Logger error:%v\n", err)
	}
	log.Printf("init zap Logger for gin mode: %s\n", gin.Mode())

	defer func(Logger *zap.Logger) {
		err := Logger.Sync()
		if err != nil {
			log.Printf("defer Logger error:%v\n", err)
		}
	}(Logger)

	SugarLogger = Logger.Sugar()
	defer func(SugarLogger *zap.SugaredLogger) {
		err := SugarLogger.Sync()
		if err != nil {
			log.Printf("defer Sugar Logger error:%v\n", err)
		}
	}(SugarLogger)
}
