package main

import (
	"github.com/WuShaoQiang/crawler/boss/controller"
	"github.com/WuShaoQiang/crawler/boss/model"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	db := model.ConnectToDB()
	defer db.Close()
	model.SetDB(db)

	controller.Register()
	controller.StartUp()
}
