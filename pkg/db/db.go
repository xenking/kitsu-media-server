package db

import (
	"fmt"

	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/xenking/kitsu-media-server/pkg/config"
	"github.com/xenking/kitsu-media-server/pkg/model"
)

func New(config *config.Config) *gorm.DB {
	cfgStr := "host=" + config.Database.Host + " port=" + config.Database.Port +
		" user=" + config.Database.Username + " dbname=" + config.Database.Name +
		" password=" + config.Database.Password + " sslmode=disable"
	db, err := gorm.Open("postgres", cfgStr)
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	db.DB().SetMaxIdleConns(3)
	db.LogMode(true)
	return db
}

func TestDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./../kitsu_media_test.db")
	if err != nil {
		fmt.Println("storage err: ", err)
	}
	db.DB().SetMaxIdleConns(3)
	db.LogMode(false)
	return db
}

func DropTestDB() error {
	if err := os.Remove("./../kitsu_media_test.db"); err != nil {
		return err
	}
	return nil
}

//TODO: err check
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.User{},
		&model.Follow{},
		&model.Article{},
		&model.Media{},
		&model.Comment{},
		&model.Tag{},
	)
}
