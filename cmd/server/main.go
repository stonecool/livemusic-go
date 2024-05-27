package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/livemusic-go/internal/router"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}

func main() {
	gin.SetMode("debug")

	engine := gin.New()
	if err := router.InitRouter(engine); err != nil {
		log.Fatal(err)
	}

	maxHeaderBytes := 1 << 20
	server := &http.Server{
		Addr:           ":8080",
		Handler:        engine,
		ReadTimeout:    300,
		WriteTimeout:   time.Second,
		MaxHeaderBytes: maxHeaderBytes,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("init server error:%v\n", err)
		return
	}

}
