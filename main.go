package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/chenjiandongx/go-echarts/charts"

	"github.com/PuerkitoBio/goquery"
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
	headers  = []string{"Host", "User-Agent", "Accept", "Accept-Language", "Accept-Encoding", "Referer", "Connection", "Cookie", "Upgrade-Insecure-Requests", "TE"}
	keywords = []string{"golang"}
	path     = "/home/shelljo/go/src/github.com/WuShaoQiang/crawler/boss/"
	// jobs     []Job
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
	url := "https://www.zhipin.com/job_detail/?query=%s&city=101010100&industry=&position="
	for index, keyword := range keywords {
		log.Printf("Running on %v\n", index)
		req, err := http.NewRequest("GET", fmt.Sprintf(url, keyword), nil)
		if err != nil {
			log.Fatalf("NewRequest Error: %s", err)
		}

		for _, header := range headers {
			req.Header.Add(header, viper.GetString(header))
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalf("Do Error %s\n", err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("ReadAll Error %s\n", err)
		}

		// fmt.Println(string(body))
		file, err := os.Create(path + keyword + ".html")
		if err != nil {
			log.Fatalf("Create File Error %s\n", err)
		}

		defer file.Close()

		_, err = file.Write(body)
		if err != nil {
			log.Fatalf("Write Error %s\n", err)
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatalf("go query Error %s\n", err)
		}

		num := doc.Find("ul>li>div.job-primary").Size()
		fmt.Println(num)

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

		fmt.Println(jobs)
		count := make([]int, 0)
		nameItems := make([]string, 0)
		countMap := make(map[string]int)
		for _, single := range jobs {
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

		show(nameItems, count)

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
