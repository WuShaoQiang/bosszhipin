package model

import (
	"log"

	"github.com/WuShaoQiang/crawler/boss/config"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

// SetDB func
func SetDB(database *gorm.DB) {
	db = database
}

// ConnectToDB func
func ConnectToDB() *gorm.DB {
	connectingStr := config.GetMysqlConnectingString()
	log.Println("Connet to db...")
	db, err := gorm.Open("mysql", connectingStr)
	if err != nil {
		log.Fatalf("Failed to connect database : %s", err)
	}
	db.SingularTable(true)
	return db
}
