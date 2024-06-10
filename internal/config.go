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
var AccountConfigMap map[string]AccountConfig

type conf struct {
	Server           server
	Database         database
	Redis            redis
	AccountConfigMap map[string]AccountConfig
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

type AccountConfig struct {
	Name           string
	Type           uint8
	CheckLoginURL  string
	LoginURL       string
	HttpMethod     string
	QRCodeSelector string
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
	AccountConfigMap = config.AccountConfigMap
}
