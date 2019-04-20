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
	usefulIPAddress string
)

func doRequest(url string) *http.Response {
	c := getNewClient()
	req := getNewGETRequest(url)
	resp, err := c.Do(req)
	if err != nil {
		if resp != nil {
			logger.Warnf("Get Error %s\n Code : %v\n Time : %v", err, resp.StatusCode, loopCounter)
		} else {
			logger.Warnf("Get Error %s\n Time : %v", err, loopCounter)
		}
		loadUsefulIPAddress()
		return doRequest(url)
	}
	return resp
}

// getNextPage return whether it has a next page,and the path of next page.
// it also will store jobs in mysql
func getNextPage(page string, keyword string, pageChannel chan string) {
	defer wg.Done()
	var nextPage string
	currentURL := fmt.Sprintf(indexURL, page)
	logger.Infoln("getting page ", currentURL)
	loadUsefulIPAddress()
	resp := doRequest(currentURL)

	logger.Infoln("Get response")

	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	//犯过一个错误，把Client的超时设置为10秒，doc在读取Body的时候刚好超时了，导致发生错误
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.Infoln(resp.StatusCode, err)
		panic("go query Error ")
	}

	if doc.Find("a[class=next]").Size() == 1 {
		doc.Find("a[class=next]").Each(func(i int, s *goquery.Selection) {
			nextPage = s.Get(i).Attr[0].Val + "&" + s.Get(i).Attr[1].Val
		})
		logger.Infoln("find another page, pass it to channel")
		pageChannel <- nextPage
	} else {
		close(pageChannel)
	}

	storeJobs(doc, keyword)
	fmt.Println("store data to mysql")
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
			jobs[i].Work = work[0]
		} else {
			jobs[i].Work = "经验不限"
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
		wg.Add(1)
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
	defer wg.Done()
	resp := doRequest(fmt.Sprintf(indexURL, page))
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

func loadUsefulIPAddress() {
	for {
		if getProxyIPAddress() {
			break
		}
	}
}

func getProxyIPAddress() bool {
	ip := GetIP()
	if ip == "" {
		return false
	}
	usefulIPAddress = ip
	return true
}

func getNewClient() *http.Client {
	logger.Infoln("Using ", usefulIPAddress, " to do request")
	urli := url.URL{}
	urlproxy, _ := urli.Parse("http://" + usefulIPAddress)
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlproxy),
		},
		// Timeout: 10 * time.Second,
	}
	return client
}

func getNewGETRequest(currentURL string) *http.Request {
	req, err := http.NewRequest("GET", currentURL, nil)
	if err != nil {
		logger.Debugln("Create Request Error : ", err)
	}

	req.Header.Set("User-Agent", GetUserAgent())
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "identity")

	return req
}

func salaryAverage(jobs []Job) int {
	n := len(jobs)
	if n == 0 {
		return 0
	}
	total := 0
	for _, job := range jobs {
		total += job.Salary
	}
	return (total / n)
}

func clearAllData(tableName string) error {
	if err := deleteTable(tableName); err != nil {
		return err
	}
	if err := createJobTable(); err != nil {
		return err
	}
	return nil
}
