package model

import (
	"fmt"
	"log"
	"sync"
)

var (
	wg sync.WaitGroup
)

// CrawlerGo start crawler and store data in mysql
// Before crawler start it will detele all old data
func CrawlerGo(keywords []string, refresh bool) bool {
	urlsEncoded := keywordEncode(keywords)
	for index, keyword := range keywords {
		if isKeywordExist(keyword) {
			if !refresh {
				continue
			}
			deleteByKeyword(keyword)
		}
		wg.Add(1)
		go crawlerGoSingleKeyword(keyword, urlsEncoded[index])
	}
	wg.Wait()
	log.Println("Crawler finished!")
	return true
}

func crawlerGoSingleKeyword(keyword, urlEncoded string) {
	defer wg.Done()
	indexPage := "/c100010000/?query=%s&page=1&ka=page-1"
	currentPage := fmt.Sprintf(indexPage, urlEncoded)
	log.Printf("Collecting on %v\n", keyword)
	for {
		next, nextPage := getNextPage(currentPage, keyword)
		if !next {
			log.Println("The end of ", keyword, " page")
			break
		} else {
			currentPage = nextPage
		}
	}
}

func isKeywordExist(keyword string) bool {
	var num int
	if err := db.Model(&Job{}).Where("keyword = ?", keyword).Count(&num).Error; err != nil {
		log.Fatalln("isKeywordExist Error : ", err)
	}
	if num > 0 {
		return true
	}
	if num < 0 {
		log.Fatalln("num can't be negative")
	}

	return false
}

func deleteByKeyword(keyword string) {
	log.Println("deleting ", keyword, " in database")
	if err := db.Table("job").Delete(&Job{}, "keyword = ?", keyword); err != nil {
		log.Fatalln("deleteByKeyword Error : ", err)
	}
}

// func loadDataToMysql() error {
// 	clearTableData("job")
// 	for _, job := range allJobs {
// 		if err := job.AddJob(); err != nil {
// 			log.Println("loadDataToMysql Error ", err)
// 			return err
// 		}
// 	}
// 	return nil
// }
