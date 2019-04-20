package model

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
	Work          string `gorm:"type:varchar(8)"`
	Education     string `gorm:"type:varchar(16)"`
	Detail        string `gorm:"type:varchar(64)"`
	Benefit       string `gorm:"type:varchar(256)"`
}

var (
	// allJobs     []Job
	indexURL  = "https://www.zhipin.com%s"
	work      = []string{"经验不限", "1-3", "3-5", "5-10"}
	education = []string{"大专", "本科", "硕士"}
	city      = []string{"北京", "上海", "广州", "深圳", "杭州", "西安", "武汉", "成都", "南京", "厦门"}
)

// AddJob add one job to database
func (job *Job) AddJob() error {
	if err := db.Create(&job).Error; err != nil {
		logger.Debugln("AddJob Error ", err)
	}
	return nil
}

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
				logger.Debugln("BarDataCityJobNum Error : ", err)
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

func BarDataSalaryWork(keywords []string) (nameItems []string, count [][]int) {
	nameItems = work
	count = make([][]int, len(keywords))

	for idx1, keyword := range keywords {
		var jobs []Job
		for _, nameItem := range nameItems {
			db.Model(&Job{}).Where("keyword = ?", keyword).Where("work = ?", nameItem).Find(&jobs)
			if count[idx1] == nil {
				count[idx1] = make([]int, 0)
			}
			average := salaryAverage(jobs)
			count[idx1] = append(count[idx1], average)
		}
	}
	return
}

func BarDataSalaryEducation(keywords []string) (nameItems []string, count [][]int) {
	nameItems = education
	count = make([][]int, len(keywords))

	for idx1, keyword := range keywords {
		var jobs []Job
		for _, nameItem := range nameItems {
			db.Model(&Job{}).Where("keyword = ?", keyword).Where("education = ?", nameItem).Find(&jobs)
			if count[idx1] == nil {
				count[idx1] = make([]int, 0)
			}
			average := salaryAverage(jobs)
			count[idx1] = append(count[idx1], average)
		}
	}
	return
}

func BarDataSalaryCity(keywords []string) (nameItems []string, count [][]int) {
	nameItems = city
	count = make([][]int, len(keywords))

	for idx1, keyword := range keywords {
		var jobs []Job
		for _, nameItem := range nameItems {
			db.Model(&Job{}).Where("keyword = ?", keyword).Where("city = ?", nameItem).Find(&jobs)
			if count[idx1] == nil {
				count[idx1] = make([]int, 0)
			}
			average := salaryAverage(jobs)
			count[idx1] = append(count[idx1], average)
		}
	}
	return
}

/*-------------------------Map-------------------------------*/

// MapDataProvinceJobNum return how many jobs are there in each province
func MapDataProvinceJobNum(keyword string) map[string]float32 {
	mapData := make(map[string]float32)
	var num int
	for city, province := range provinceMap {
		err := db.Model(&Job{}).Where("city = ?", city).Where("keyword = ?", keyword).Count(&num).Error
		if err != nil {
			logger.Debugln("MapData Error : ", err)
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
