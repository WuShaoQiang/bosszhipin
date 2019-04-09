package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chenjiandongx/go-echarts/charts"
	"github.com/spf13/viper"
)

type Job struct {
	// ID        int
	Salary    int
	Location  string
	Wrok      string
	Education string
}

var (
	headers  = []string{"Host", "User-Agent", "Accept", "Accept-Language", "Accept-Encoding", "Connection"}
	keywords = []string{"golang", "python"}
	path     = "/home/shelljo/go/src/github.com/WuShaoQiang/crawler/boss/"
	allJobs  []Job
	url      = "https://www.zhipin.com%s"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	readHeader()
}

func readHeader() {
	viper.SetConfigName("header")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("readHeader Error %s\n", err)
	}
}

func main() {

	index_page := "/c101010100/?query=%s&page=1&ka=page-1"
	for index, keyword := range keywords {
		currentPage := fmt.Sprintf(index_page, keyword)
		log.Printf("Running on %v\n", index)
		for {
			next, next_page := getNextPage(currentPage)
			if !next {
				break
			} else {
				currentPage = next_page
			}
		}

		fmt.Println(allJobs)
		show(counter())

	}
}

func show(nameItems []string, count []int) {
	// nameItems := []string{"北京"}
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.TitleOpts{Title: "Golang收入情况"})
	bar.AddXAxis(nameItems).AddYAxis("北京", count)
	f, err := os.Create(path + "bar.html")
	if err != nil {
		log.Printf("show Error %s\n", err)
	}
	bar.Render(f)
}

func getNextPage(page string) (bool, string) {
	var next_page string
	resp, err := http.Get(fmt.Sprintf(url, page))
	if err != nil {
		log.Fatalf("Get Error %s\n", err)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("go query Error %s\n", err)
	}

	loadData(doc)
	// <a href="javascript:;" ka="page-next" class="next disabled"></a>
	// <a href="/c101010100/?query=golang&amp;page=10" ka="page-next" class="next"></a>
	// fmt.Println(doc.Find("a[class=next]").Html())
	if doc.Find("a[class=next]").Size() == 1 {
		doc.Find("a[class=next]").Each(func(i int, s *goquery.Selection) {

			// fmt.Println(i, s.Get(i).Attr[0].Val+"&"+s.Get(i).Attr[1].Val)
			next_page = s.Get(i).Attr[0].Val + "&" + s.Get(i).Attr[1].Val
		})
		return true, next_page
	}
	return false, ""
}

func loadData(doc *goquery.Document) error {
	num := doc.Find("ul>li>div.job-primary").Size()
	// fmt.Println(num)

	jobs := make([]Job, num)

	doc.Find("ul>li>div>div>h3>a>span").Each(func(i int, s *goquery.Selection) {
		str := s.Text()
		str = strings.Replace(str, "k", "", -1)
		strs := strings.Split(str, "-")
		// fmt.Println(strs)
		num1, _ := strconv.Atoi(strs[0])
		num2, _ := strconv.Atoi(strs[1])

		jobs[i].Salary = (num1 + num2) / 2
	})

	doc.Find("ul>li>div>div.info-primary>p").Each(func(i int, s *goquery.Selection) {
		// fmt.Println(i, s.Text())
		str := s.Text()
		// reg := regexp.MustCompile(`[\p{Han}]+`)
		// fmt.Printf("%q\n", reg.FindAllString(str, -1))
		// jobs[i].Location = s.Text()

		reg1 := regexp.MustCompile(`\d-\d+`)
		// fmt.Printf("%q\n", reg1.FindAllString(str, -1))
		work := reg1.FindAllString(str, -1)
		if len(work) > 1 {
			log.Println("workExperience too long")
		} else if len(work) == 1 {
			jobs[i].Wrok = work[0]
		} else {
			jobs[i].Wrok = "经验不限"
		}

		// reg2 := regexp.MustCompile(`([\p{Han}]+\s)`)
		// fmt.Printf("%q\n", reg2.FindAllString(str, -1))

		// nian := fmt.Sprintf("%x", "年")
		// fmt.Println(nian)
		reg3 := regexp.MustCompile(`[^\w\s\-]+`)
		// fmt.Printf("%q\n", reg3.FindAllString(str, -1))
		temp := reg3.FindAllString(str, -1)
		location := temp[0]
		// fmt.Println(location)
		jobs[i].Location = location
	})

	allJobs = append(allJobs, jobs...)

	return nil
}

func counter() (nameItems []string, count []int) {
	// count := make([]int, 0)
	// nameItems := make([]string, 0)
	countMap := make(map[string]int)
	for _, single := range allJobs {
		if _, ok := countMap[strconv.Itoa(single.Salary)]; !ok {
			countMap[strconv.Itoa(single.Salary)] = 1
			nameItems = append(nameItems, strconv.Itoa(single.Salary))
		} else {
			countMap[strconv.Itoa(single.Salary)]++
		}
	}

	for _, single := range nameItems {
		count = append(count, countMap[single])
	}

	return
}
