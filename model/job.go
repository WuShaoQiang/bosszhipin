package model

import (
	"fmt"
	"log"
	"strconv"
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

func MapData() map[string]float32 {
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
