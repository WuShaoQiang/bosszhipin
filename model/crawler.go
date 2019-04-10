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
func CrawlerGo(keywords, urlsEncoded []string) {
	log.Println("Cleaning Database...")
	if err := clearTableData("job"); err != nil {
		log.Fatalln("CrawlerGo Error : ", err)
	}
	for index, keyword := range keywords {
		wg.Add(1)
		go crawlerGoSingleKeyword(keyword, urlsEncoded[index])
	}
	wg.Wait()
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
