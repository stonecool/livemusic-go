package internal

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

var config conf

var Server = &config.Server
var Database = &config.Database
var Redis = &config.Redis
var CrawlAccountMap map[string]CrawlAccount

type conf struct {
	Server          server
	Database        database
	Redis           redis
	CrawlAccountMap map[string]CrawlAccount
}

type server struct {
	RunMode      string
	HttpPort     uint
	ReadTimeout  uint
	WriteTimeout uint
}

type database struct {
	Type         string
	User         string
	Password     string
	Host         string
	DatabaseName string
	TablePrefix  string
	Charset      string
	ParseTime    bool
	Loc          string
}

type redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout int
}

type CrawlAccount struct {
	Name          string
	Type          uint8
	LoginURL      string
	CheckLoginURL string
}

func init() {
	_, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	_, err = toml.DecodeFile("conf/conf.toml", &config)
	if err != nil {
		fmt.Println(err)
	}

	log.Println(config)
	CrawlAccountMap = config.CrawlAccountMap
}
