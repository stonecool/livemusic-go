package model

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

var db *gorm.DB

type Model struct {
	ID        int `gorm:"primary key"`
	CreatedAt int
	UpdatedAt int
	DeletedAt int
}

// init initializes the database instance
func init() {
	var dialector gorm.Dialector
	if internal.Database.Type == "mysql" {
		dialector = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
			internal.Database.User,
			internal.Database.Password,
			internal.Database.Host,
			internal.Database.DatabaseName,
			internal.Database.Charset,
			internal.Database.ParseTime,
			internal.Database.Loc))
	}

	var err error
	db, err = gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   internal.Database.TablePrefix,
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
