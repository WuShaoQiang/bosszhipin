package model

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// getNextPage return whether it has a next page,and the path of next page.
func getNextPage(page string) (bool, string) {
	var nextPage string
	resp, err := http.Get(fmt.Sprintf(url, page))
	if err != nil {
		log.Fatalf("Get Error %s\n", err)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("go query Error %s\n", err)
	}

	loadDataToVar(doc)
	if doc.Find("a[class=next]").Size() == 1 {
		doc.Find("a[class=next]").Each(func(i int, s *goquery.Selection) {
			nextPage = s.Get(i).Attr[0].Val + "&" + s.Get(i).Attr[1].Val
		})
		return true, nextPage
	}
	return false, ""
}

func isExist(nameItems []string, item string) bool {
	for _, single := range nameItems {
		if single == item {
			return true
		}
	}
	return false
}
