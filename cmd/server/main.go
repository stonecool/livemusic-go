package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stonecool/1701livehouse-server/internal/router"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}

func main() {
	gin.SetMode("debug")

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

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

	//log.Printf("[info] start http server listening %s", endPoint)
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
