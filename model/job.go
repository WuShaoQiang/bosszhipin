package model

import (
	"fmt"
	"log"
	"strconv"
)

// Job struct
type Job struct {
	ID            int    `gorm:"primary_key"`
	Name          string `gorm:"type:varchar(32)"`
	Salary        int    `gorm:"type:int(16)"`
	Company       string `gorm:"type:varchar(64)"`
	Province      string `gorm:"type:varchar(16)"`
	City          string `gorm:"type:varchar(16)"`
	District      string `gorm:"type:varchar(16)"`
	DetailAddress string `gorm:"type:varchar(64)"`
	Wrok          string `gorm:"type:varchar(8)"`
	Education     string `gorm:"type:varchar(16)"`
	Detail        string `gorm:"type:varchar(64)"`
	Benefit       string `gorm:"type:varchar(256)"`
}

var (
	allJobs     []Job
	url         = "https://www.zhipin.com%s"
	keywords    = []string{"golang"}
	provinceMap = map[string]string{
		"深圳":  "广东",
		"广州":  "广东",
		"珠海":  "广东",
		"佛山":  "广东",
		"厦门":  "福建",
		"福州":  "福建",
		"北京":  "北京",
		"杭州":  "浙江",
		"温州":  "浙江",
		"宁波":  "浙江",
		"天津":  "天津",
		"上海":  "上海",
		"重庆":  "重庆",
		"西安":  "陕西",
		"南京":  "江苏",
		"武汉":  "湖北",
		"石家庄": "河北",
		"合肥":  "安徽",
		"长沙":  "湖南",
		"郑州":  "河南",
		"济南":  "山东",
		"香港":  "香港",
		"澳门":  "澳门",
		"台北":  "台湾",
		"南昌":  "江西",
		"南宁":  "广西",
		"贵阳":  "贵州",
		"成都":  "四川",
		"沈阳":  "辽宁",
		"大连":  "辽宁",
		"太原":  "山西",
		"兰州":  "甘肃",
		"昆明":  "云南",
		"银川":  "宁夏",
		"西宁":  "青海",
		"南通":  "江苏",
		"德州":  "山东",
		"柳州":  "广西",
	}
)

// // AddJob add one job to database
// func AddJob(salary int, location, work, education string) error {
// 	job := Job{Salary: salary, City: location, Wrok: work, Education: education}
// 	if err := db.Create(&job).Error; err != nil {
// 		log.Println("AddJob Error ", err)
// 		return err
// 	}
// 	return nil
// }

// func loadDataToMysql() error {
// 	for _, job := range allJobs {
// 		if err := AddJob(job.Salary, job.City, job.Wrok, job.Education); err != nil {
// 			log.Println("loadDataToMysql Error ", err)
// 			return err
// 		}
// 	}
// 	return nil
// }

// AddJob add one job to database
func (job *Job) AddJob() error {
	if err := db.Create(&job).Error; err != nil {
		log.Println("AddJob Error ", err)
		return err
	}
	return nil
}

func loadDataToMysql() error {
	clearTableData("job")
	for _, job := range allJobs {
		if err := job.AddJob(); err != nil {
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
		if _, exist := cityCountMap[single.City]; !exist {
			cities = append(cities, single.City)

			city2NumMap[single.City] = counter
			countMap = append(countMap, map[string]int{})
			fmt.Println(counter)
			cityCountMap[single.City] = countMap[counter]
			counter++
		} else {
			if _, exist := countMap[city2NumMap[single.City]][strconv.Itoa(single.Salary)]; !exist {
				countMap[city2NumMap[single.City]][strconv.Itoa(single.Salary)] = 1
				if !isExist(nameItems, strconv.Itoa(single.Salary)) {
					nameItems = append(nameItems, strconv.Itoa(single.Salary))
				}
			} else {
				countMap[city2NumMap[single.City]][strconv.Itoa(single.Salary)]++
			}
		}

	}

	return
}

func MapData() map[string]float32 {
	mapData := make(map[string]float32)
	for _, job := range allJobs {
		if _, exist := mapData[job.City]; !exist {
			mapData[job.City] = 1
		} else {
			mapData[job.City]++
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
		// getNextPage(currentPage)
		loadDataToMysql()
		// fmt.Println(allJobs)
		// show(counter())

		// getDetailPage("/job_detail/0fb6f3bb0f577d031X192tq4ElQ~.html")

	}
}
