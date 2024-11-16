package database

import (
	"fmt"
	"github.com/stonecool/livemusic-go/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

var DB *gorm.DB

type BaseModel struct {
	ID        int `gorm:"primary key"`
	CreatedAt int
	UpdatedAt int
	DeletedAt int
}

// init initializes the database instance
// TODO func Initialize() {?
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

func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
