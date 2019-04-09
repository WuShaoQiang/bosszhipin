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

// Job struct
type Job struct {
	// ID        int
	Salary    int
	Location  string
	Wrok      string
	Education string
}

var (
	headers  = []string{"Host", "User-Agent", "Accept", "Accept-Language", "Accept-Encoding", "Connection"}
	keywords = []string{"golang"}
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
	indexPage := "/c100010000/?query=%s&page=1&ka=page-1"
	for index, keyword := range keywords {
		currentPage := fmt.Sprintf(indexPage, keyword)
		log.Printf("Running on %v\n", index)
		for {
			next, nextPage := getNextPage(currentPage)
			if !next {
				break
			} else {
				currentPage = nextPage
			}
		}

		// fmt.Println(allJobs)
		// show(counter())

	}
}

func show(nameItems []string, cityCountMap map[string]map[string]int, cities []string) {
	bar := charts.NewBar()
	bar.SetGlobalOptions(charts.TitleOpts{Title: "Golang收入情况"})
	tmp := bar.AddXAxis(nameItems)
	for _, city := range cities {
		tmpArray := make([]int, 0)
		for _, nameItem := range nameItems {
			if num, exist := cityCountMap[city][nameItem]; exist {
				tmpArray = append(tmpArray, num)
			} else {
				tmpArray = append(tmpArray, 0)
			}
		}
		tmp.AddYAxis(city, tmpArray)
	}
	f, err := os.Create(path + "bar.html")
	if err != nil {
		log.Printf("show Error %s\n", err)
	}
	bar.Render(f)
}

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

	loadData(doc)
	// <a href="javascript:;" ka="page-next" class="next disabled"></a>
	// <a href="/c101010100/?query=golang&amp;page=10" ka="page-next" class="next"></a>
	// fmt.Println(doc.Find("a[class=next]").Html())
	if doc.Find("a[class=next]").Size() == 1 {
		doc.Find("a[class=next]").Each(func(i int, s *goquery.Selection) {

			// fmt.Println(i, s.Get(i).Attr[0].Val+"&"+s.Get(i).Attr[1].Val)
			nextPage = s.Get(i).Attr[0].Val + "&" + s.Get(i).Attr[1].Val
		})
		return true, nextPage
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

func counter() (nameItems []string, cityCountMap map[string]map[string]int, cities []string) {
	cityCountMap = map[string]map[string]int{}
	countMap := make([]map[string]int, 0)
	// for index := range countMap {
	// 	countMap[index] = make(map[string]int)
	// }
	// fmt.Println(countMap[2])
	city2NumMap := map[string]int{}
	counter := 0
	for _, single := range allJobs {
		if _, exist := cityCountMap[single.Location]; !exist {
			cities = append(cities, single.Location)

			city2NumMap[single.Location] = counter
			countMap = append(countMap, map[string]int{})
			fmt.Println(counter)
			cityCountMap[single.Location] = countMap[counter]
			counter++
		} else {
			if _, exist := countMap[city2NumMap[single.Location]][strconv.Itoa(single.Salary)]; !exist {
				countMap[city2NumMap[single.Location]][strconv.Itoa(single.Salary)] = 1
				if !isExist(nameItems, strconv.Itoa(single.Salary)) {
					nameItems = append(nameItems, strconv.Itoa(single.Salary))
				}
			} else {
				countMap[city2NumMap[single.Location]][strconv.Itoa(single.Salary)]++
			}
		}

	}

	return
}

func isExist(nameItems []string, item string) bool {
	for _, single := range nameItems {
		if single == item {
			return true
		}
	}
	return false
}

func countryMap() map[string]float32 {
	mapData := make(map[string]float32)
	for _, job := range allJobs {
		if _, exist := mapData[job.Location]; !exist {
			mapData[job.Location] = 1
		} else {
			mapData[job.Location]++
		}
	}
	return mapData
}

func mapVisualMap(mapData map[string]float32) *charts.Map {
	mc := charts.NewMap("china")
	mc.SetGlobalOptions(
		charts.TitleOpts{Title: "Map-设置 VisualMap"},
		charts.VisualMapOpts{Calculable: true},
	)
	mc.Add("map", mapData)
	return mc
}
