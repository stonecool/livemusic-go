package model

import (
	"fmt"
	"github.com/stonecool/1701livehouse-server/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

var db *gorm.DB

type Model struct {
	ID        int `gorm:"primary key" json:"id"`
	CreatedAt int `json:"created_at"`
	UpdatedAt int `json:"updated_at"`
	DeletedAt int `json:"deleted_at"`
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
	db, err = gorm.Open(dialector, &gorm.Config{
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
	//defer db()
}
