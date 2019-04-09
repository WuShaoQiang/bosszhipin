package main

import (
	"log"

	"github.com/WuShaoQiang/crawler/boss/model"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	log.Println("DB Init ...")
	db := model.ConnectToDB()
	defer db.Close()
	model.SetDB(db)

	db.DropTableIfExists(model.Job{})
	db.CreateTable(model.Job{})

	model.AddJob(20, "北京", "3-5", "")

}
