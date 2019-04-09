package model

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Job struct
type Job struct {
	ID        int    `gorm:"primary_key"`
	Salary    int    `gorm:"type:int(16)"`
	Location  string `gorm:"type:varchar(32)"`
	Wrok      string `gorm:"type:varchar(8)"`
	Education string `gorm:"type:varchar(32)"`
}

var (
	allJobs  []Job
	url      = "https://www.zhipin.com%s"
	keywords = []string{"golang"}
)

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
	doc.Find("ul>li>div>div.info-primary>p").Each(func(i int, s *goquery.Selection) {
		str := s.Text()
		reg1 := regexp.MustCompile(`\d-\d+`)
		work := reg1.FindAllString(str, -1)
		if len(work) > 1 {
			log.Println("workExperience too long")
		} else if len(work) == 1 {
			jobs[i].Wrok = work[0]
		} else {
			jobs[i].Wrok = "经验不限"
		}

		reg3 := regexp.MustCompile(`[^\w\s\-]+`)
		temp := reg3.FindAllString(str, -1)
		location := temp[0]
		jobs[i].Location = location
	})

	// Add to allJobs
	allJobs = append(allJobs, jobs...)

	return nil
}

// AddJob add one job to database
func AddJob(salary int, location, work, education string) error {
	job := Job{Salary: salary, Location: location, Wrok: work, Education: education}
	if err := db.Create(&job).Error; err != nil {
		log.Println("AddJob Error ", err)
		return err
	}
	return nil
}

func loadDataToMysql() error {
	for _, job := range allJobs {
		if err := AddJob(job.Salary, job.Location, job.Wrok, job.Education); err != nil {
			log.Println("loadDataToMysql Error ", err)
			return err
		}
	}
	return nil
}

func Counter() (nameItems []string, cityCountMap map[string]map[string]int, cities []string) {
	cityCountMap = map[string]map[string]int{}
	countMap := make([]map[string]int, 0)
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

func Process() {
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
		loadDataToMysql()

		// fmt.Println(allJobs)
		// show(counter())

	}
}
