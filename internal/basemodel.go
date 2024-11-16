package internal

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

var DB *gorm.DB

type RawModel struct {
	ID        int `gorm:"primary key"`
	CreatedAt int
	UpdatedAt int
	DeletedAt int
}

// init initializes the database instance
func init() {
	var dialector gorm.Dialector
	if config.Database.Type == "mysql" {
		dialector = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
			config.Database.User,
			config.Database.Password,
			config.Database.Host,
			config.Database.DatabaseName,
			config.Database.Charset,
			config.Database.ParseTime,
			config.Database.Loc))
	}

	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   config.Database.TablePrefix,
			//NameReplacer:
			//NoLowerCase:
		},
	})
	if err != nil {
		panic(err)
	} else {
		log.Println("database connect finish!")
	}

}

// TODO
func closeDB() {
	//defer DB()
}
