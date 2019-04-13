package model

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

var (
	loopCounter = 0
	// c           http.Client
)

// getNextPage return whether it has a next page,and the path of next page.
// it also will store jobs in mysql
func getNextPage(page string, keyword string) (bool, string) {
	var nextPage string
	currentURL := fmt.Sprintf(indexURL, page)
	fmt.Println("getting page ", currentURL)
	// resp, err := http.Get(currentURL)
	c, currentIP := getNewClient()
	req := getNewGETRequest(currentURL)

	resp, err := c.Do(req)
	if err != nil || resp.StatusCode != 200 {
		logger.Warnf("Get Error %s\n Code : %v\n Time : %v", err, resp.StatusCode, loopCounter)
		DeleteIP(currentIP)
		if loopCounter > 3 {
			logger.Debug("Made request more than three times, stopping the program")
		}
		loopCounter++
		return getNextPage(page, keyword)
	}

	//reset loopCounter if store data successfully
	loopCounter = 0

	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.Debugln("go query Error ", err)
	}

	storeJobs(doc, keyword)

	fmt.Println("store data to mysql")
	if doc.Find("a[class=next]").Size() == 1 {
		doc.Find("a[class=next]").Each(func(i int, s *goquery.Selection) {
			nextPage = s.Get(i).Attr[0].Val + "&" + s.Get(i).Attr[1].Val
		})
		return true, nextPage
	}
	return false, ""
}

func storeJobs(doc *goquery.Document, keyword string) error {
	num := doc.Find("ul>li>div.job-primary").Size()
	jobs := make([]Job, num)

	//Find company
	doc.Find("ul>li>div>div.info-company>div>h3>a").Each(func(i int, s *goquery.Selection) {
		jobs[i].Company = s.Text()
	})

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
	// Find city
	// Find district
	// Find province
	// Find education
	doc.Find("ul>li>div>div.info-primary>p").Each(func(i int, s *goquery.Selection) {
		str := s.Text()
		r := []rune(str)
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

		strSplit := strings.Split(str, " ")

		switch len(strSplit) {
		case 2:
			jobs[i].City = strSplit[0]
		case 3:
			jobs[i].City = strSplit[0]
			jobs[i].District = strSplit[1]

		default:
			log.Println("Something wrong in analysing city and district")
		}

		if province, exist := provinceMap[jobs[i].City]; !exist {
			log.Println(jobs[i].City, " doesn't exist in provinceMap")
			jobs[i].Province = "未知"
		} else {
			jobs[i].Province = province
		}

	})

	//Find detail
	doc.Find("ul>li>div.job-primary>div>h3>a").Each(func(i int, s *goquery.Selection) {
		detail := s.Get(0).Attr[0].Val
		jobs[i].Detail = detail
		go jobs[i].getDetailPage(detail)
	})

	//Find name
	doc.Find("ul>li>div.job-primary>div>h3>a>div.job-title").Each(func(i int, s *goquery.Selection) {
		jobs[i].Name = s.Text()
	})

	for _, job := range jobs {
		job.Keyword = keyword
		if err := job.AddJob(); err != nil {
			return err
		}
	}

	return nil
}

// getDetailPage is for job's benefit and detailAddress
func (job *Job) getDetailPage(page string) {
	c, currentIP := getNewClient()
	req := getNewGETRequest(fmt.Sprintf(indexURL, page))
	resp, err := c.Do(req)
	if err != nil || resp.StatusCode != 200 {
		logger.Warnf("Get Error %s\n Code : %v\n Time : %v", err, resp.StatusCode, loopCounter)
		DeleteIP(currentIP)
		if loopCounter > 3 {
			logger.Debug("Made request more than three times, stopping the program")
		}
		loopCounter++
		job.getDetailPage(page)
	}

	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.Debugln("go query Error ", err)
	}

	job.Benefit, job.DetailAddress = getDetailData(doc)
}

// getDatailData return benefit and detailAddress
func getDetailData(doc *goquery.Document) (benefit string, detailAddress string) {
	//Find benefit
	var single []string
	doc.Find("div.info-primary>div.tag-container>div.job-tags>span").Each(func(i int, s *goquery.Selection) {
		single = append(single, s.Text())

	})
	benefit = strings.Join(single, ",")

	doc.Find("div.job-detail>div.detail-content>div.job-sec>div.job-location>div.location-address").Each(func(i int, s *goquery.Selection) {
		detailAddress = s.Text()
	})

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

func deleteTable(tableName string) error {
	return errors.Wrap(db.DropTableIfExists(tableName).Error, "deteleTable Error")
}

func createJobTable() error {
	return errors.Wrap(db.CreateTable(&Job{}).Error, "createTable Error")
}

func keywordEncode(keywords []string) (urlsEncoded []string) {
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

func getNewClient() (http.Client, string) {
	newIP := GetIP()
	logger.Infoln("Using ", newIP, " to do request")
	urli := url.URL{}
	urlproxy, _ := urli.Parse("//" + newIP)
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlproxy),
		},
	}
	return client, newIP
}

func getNewGETRequest(currentURL string) *http.Request {
	req, err := http.NewRequest("GET", currentURL, nil)
	if err != nil {
		logger.Debugln("Create Request Error : ", err)
	}

	req.Header.Set("User-Agent", GetUserAgent())
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	return req
}
