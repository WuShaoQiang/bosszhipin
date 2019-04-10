package model

import (
	"log"
)

// Job struct
type Job struct {
	ID            int    `gorm:"primary_key"`
	Keyword       string `gorm:"type:varchar(32)"`
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
	// allJobs     []Job
	url = "https://www.zhipin.com%s"
)

// AddJob add one job to database
func (job *Job) AddJob() error {
	if err := db.Create(&job).Error; err != nil {
		log.Println("AddJob Error ", err)
		return err
	}
	return nil
}

// func BarData() (nameItems []string, cityCountMap map[string]map[string]int, cities []string) {
// 	cityCountMap = map[string]map[string]int{}
// 	countMap := make([]map[string]int, 0)
// 	city2NumMap := map[string]int{}
// 	counter := 0
// 	for _, single := range allJobs {
// 		if _, exist := cityCountMap[single.City]; !exist {
// 			cities = append(cities, single.City)

// 			city2NumMap[single.City] = counter
// 			countMap = append(countMap, map[string]int{})
// 			fmt.Println(counter)
// 			cityCountMap[single.City] = countMap[counter]
// 			counter++
// 		} else {
// 			if _, exist := countMap[city2NumMap[single.City]][strconv.Itoa(single.Salary)]; !exist {
// 				countMap[city2NumMap[single.City]][strconv.Itoa(single.Salary)] = 1
// 				if !isExist(nameItems, strconv.Itoa(single.Salary)) {
// 					nameItems = append(nameItems, strconv.Itoa(single.Salary))
// 				}
// 			} else {
// 				countMap[city2NumMap[single.City]][strconv.Itoa(single.Salary)]++
// 			}
// 		}

// 	}

// 	return
// }

// BarDataCityJobNum return how many jobs are there in each city
func BarDataCityJobNum(keywords []string) (nameItems []string, count [][]int) {
	var num int
	// 出错过一次，没有初始化count，下面的count[index]发生溢出
	count = make([][]int, len(keywords))
	// 出错过一次，初始化了数组，但是数组的map还未初始化
	cityNumMap := make([]map[string]int, len(keywords))
	for i := 0; i < len(keywords); i++ {
		cityNumMap[i] = make(map[string]int)
	}
	for index, keyword := range keywords {
		for city := range provinceMap {
			err := db.Model(&Job{}).Where("city = ?", city).Where("keyword = ?", keyword).Count(&num).Error
			if err != nil {
				log.Fatalf("BarDataCityJobNum Error : %s\n", err)
			} else {
				if num != 0 {
					cityNumMap[index][city] = num
					if !isExist(nameItems, city) {
						nameItems = append(nameItems, city)
					}
				}
			}
		}
	}
	for _, nameItem := range nameItems {
		for index := range keywords {
			if num, exist := cityNumMap[index][nameItem]; !exist {
				count[index] = append(count[index], 0)
			} else {
				count[index] = append(count[index], num)
			}
		}
	}

	return
}

// MapDataProvinceJobNum return how many jobs are there in each province
func MapDataProvinceJobNum() map[string]float32 {
	mapData := make(map[string]float32)
	var num int
	for city, province := range provinceMap {
		err := db.Model(&Job{}).Where("city = ?", city).Count(&num).Error
		if err != nil {
			log.Fatalf("MapData Error : %s\n", err)
		} else {
			if _, exist := mapData[province]; !exist {
				mapData[province] = float32(num)
			} else {
				mapData[province] += float32(num)
			}
		}
	}
	return mapData
}
