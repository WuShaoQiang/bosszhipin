package main

import (
	"log"

	"github.com/WuShaoQiang/crawler/boss/model"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	keywords = []string{"golang实习", "golang"}
)

func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func main() {
	log.Println("DB Init ...")
	db := model.ConnectToDB()
	defer db.Close()
	model.SetDB(db)
	model.CrawlerGo(keywords, true)
}

// func deleteByKeyword(keyword string) {
// 	log.Println("deleting ", keyword, " in database")
// 	if err := db.Delete(&Job{}, "keyword = ?", keyword); err != nil {
// 		log.Fatalln("deleteByKeyword Error : ", err)
// 	}

// }
