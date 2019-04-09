package model

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

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

func loadDataToVar(doc *goquery.Document) error {
	num := doc.Find("ul>li>div.job-primary").Size()
	jobs := make([]Job, num)

	// Find average salary
	doc.Find("ul>li>div>div>h3>a>span").Each(func(i int, s *goquery.Selection) {
		str := s.Text()
		str = strings.Replace(str, "k", "", -1)
		strs := strings.Split(str, "-")
		num1, _ := strconv.Atoi(strs[0])
		num2, _ := strconv.Atoi(strs[1])

		jobs[i].Salary = (num1 + num2) / 2
	})

	// Find work experience
	// Find location(city)
	// Find education
	doc.Find("ul>li>div>div.info-primary>p").Each(func(i int, s *goquery.Selection) {
		str := s.Text()
		r := []rune(str)
		fmt.Println(string(r))
		reg1 := regexp.MustCompile(`\d-\d+`)
		work := reg1.FindAllString(str, -1)
		if len(work) > 1 {
			log.Println("workExperience too long")
		} else if len(work) == 1 {
			jobs[i].Wrok = work[0]
		} else {
			jobs[i].Wrok = "经验不限"
		}

		education := string(r[len(r)-2 : len(r)])
		jobs[i].Education = education

		reg3 := regexp.MustCompile(`[^\w\s\-]+`)
		temp := reg3.FindAllString(str, -1)
		location := temp[0]
		jobs[i].City = location
	})

	doc.Find("ul>li>div.job-primary>div>h3>a").Each(func(i int, s *goquery.Selection) {
		detail := s.Get(0).Attr[0].Val
		// fmt.Println(i, nextPage)
		jobs[i].Detail = detail
	})

	// Add to allJobs
	allJobs = append(allJobs, jobs...)

	return nil
}

func isExist(nameItems []string, item string) bool {
	for _, single := range nameItems {
		if single == item {
			return true
		}
	}
	return false
}
