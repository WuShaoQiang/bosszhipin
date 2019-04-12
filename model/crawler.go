package model

import (
	"fmt"
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
	fmt.Println("Crawler finished!")
	return true
}

func crawlerGoSingleKeyword(keyword, urlEncoded string) {
	defer wg.Done()
	indexPage := "/c100010000/?query=%s&page=1&ka=page-1"
	currentPage := fmt.Sprintf(indexPage, urlEncoded)
	fmt.Printf("Collecting on %v\n", keyword)
	for {
		next, nextPage := getNextPage(currentPage, keyword)
		if !next {
			fmt.Println("The end of ", keyword, " page")
			break
		} else {
			currentPage = nextPage
		}
	}
}

func isKeywordExist(keyword string) bool {
	var num int
	if err := db.Model(&Job{}).Where("keyword = ?", keyword).Count(&num).Error; err != nil {
		logger.Debugln("isKeywordExist Error : ", err)
	}
	if num > 0 {
		return true
	}
	return false
}

func deleteByKeyword(keyword string) {
	fmt.Println("deleting ", keyword, " in database")
	if err := db.Delete(&Job{}, "keyword = ?", keyword); err != nil {
		logger.Debugln("deleteByKeyword Error : ", err)
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
