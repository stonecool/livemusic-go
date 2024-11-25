package config

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

var config conf

var Server = &config.Server
var Database = &config.Database
var Redis = &config.Redis
var AccountMap map[string]Account

type conf struct {
	Server     server
	Database   database
	Redis      redis
	AccountMap map[string]Account
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

type Account struct {
	Name     string `toml:"name"`
	Type     uint8  `toml:"type"`
	LoginURL string `toml:"login_url"`
	LastURL  string `toml:"last_url"`
}

func init() {
	// 读取配置文件路径
	configFilePath := os.Getenv("CONFIG_PATH")
	if configFilePath == "" {
		configFilePath = "../../conf/conf.toml" // 默认路径
	}

	_, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	_, err = toml.DecodeFile(configFilePath, &config)
	if err != nil {
		log.Fatal(err)
	}

	// 验证配置
	if err := validateConfig(); err != nil {
		log.Fatal(err)
	}

	AccountMap = config.AccountMap
}

func validateConfig() error {
	// 检查必需的配置项
	if config.Database.User == "" {
		return fmt.Errorf("database user cannot be empty")
	}
	// 其他验证逻辑...
	return nil
}
