package models

import (
	"banners/internal/config"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func init() {

	cfg := config.MustLoad()

	username := cfg.DataBase.Username
	password := cfg.DataBase.Password
	dbName := cfg.DataBase.Name
	dbHost := cfg.DataBase.Host
	dbType := cfg.DataBase.Type

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)

	fmt.Println(dbUri)

	conn, err := gorm.Open(dbType, dbUri)
	if err != nil {
		fmt.Println(err)
	}

	db = conn
	db.Debug().AutoMigrate(&Banner{})
}

func GetDB() *gorm.DB { return db }
