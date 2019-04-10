package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/WuShaoQiang/crawler/boss/model"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	keywords    = []string{"golang实习"}
	urlsEncoded []string
)

func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	urlsEncoded = keywordEncode()
}

func main() {
	log.Println("DB Init ...")
	db := model.ConnectToDB()
	defer db.Close()
	model.SetDB(db)
	model.CrawlerGo(keywords, urlsEncoded)
}

func keywordEncode() (urlsEncoded []string) {
	urlsEncoded = make([]string, 0)
	for _, keyword := range keywords {
		urlsEncoded = append(urlsEncoded, urlEncode(keyword))
	}
	return
}

func urlEncode(keyword string) string {
	reg := regexp.MustCompile(`[\p{Han}]+`)
	strs := reg.FindAllString(keyword, -1)
	chinese := strings.Join(strs, "")
	encoded := chineseEncode(chinese)
	encodedURL := reg.ReplaceAllString(keyword, encoded)
	return encodedURL
}

func chineseEncode(chinese string) (encoded string) {
	encoded = ""
	byteChinese := []byte(chinese)
	// return fmt.Sprintf("%x")
	for _, singleByte := range byteChinese {
		encoded = encoded + "%" + fmt.Sprintf("%x", singleByte)
	}
	return
}
