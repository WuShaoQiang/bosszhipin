package main

import (
	"fmt"
	"log"

	"github.com/WuShaoQiang/crawler/boss/model"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	keywords = []string{"golang实习", "golang", "后端开发实习生"}
)

func main() {
	log.Println("DB Init ...")
	db := model.ConnectToDB()
	defer db.Close()
	model.SetDB(db)
	fmt.Println(model.GetIP())
}
