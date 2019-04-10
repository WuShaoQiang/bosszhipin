package main

import (
	"log"

	"github.com/WuShaoQiang/crawler/boss/model"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	keywords = []string{"golang", "python", "java"}
)

func main() {
	log.Println("DB Init ...")
	db := model.ConnectToDB()
	defer db.Close()
	model.SetDB(db)
	model.CrawlerGo(keywords)
}
