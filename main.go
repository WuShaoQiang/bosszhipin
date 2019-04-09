package main

import (
	"net/http"

	"github.com/WuShaoQiang/crawler/boss/controller"
	"github.com/WuShaoQiang/crawler/boss/model"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	db := model.ConnectToDB()
	defer db.Close()
	model.SetDB(db)

	model.Process()

	controller.Register()

	http.ListenAndServe(":8080", nil)
}
