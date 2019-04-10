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
func CrawlerGo(keywords []string) {

	if err := clearTableData("job"); err != nil {
		log.Fatalln("CrawlerGo Error : ", err)
	}
	for _, keyword := range keywords {
		wg.Add(1)
		go crawlerGoSingleKeyword(keyword)
	}
	wg.Wait()
}

func crawlerGoSingleKeyword(keyword string) {
	defer wg.Done()
	indexPage := "/c100010000/?query=%s&page=1&ka=page-1"
	currentPage := fmt.Sprintf(indexPage, keyword)
	log.Printf("Collecting on %v\n", keyword)
	for {
		next, nextPage := getNextPage(currentPage, keyword)
		if !next {
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
