package vm

import (
	"github.com/WuShaoQiang/crawler/boss/model"

	"github.com/chenjiandongx/go-echarts/charts"
)

var (
// path = "/home/shelljo/go/src/github.com/WuShaoQiang/crawler/boss/"
)

// func show(nameItems []string, cityCountMap map[string]map[string]int, cities []string) {
// 	bar := charts.NewBar()
// 	bar.SetGlobalOptions(charts.TitleOpts{Title: "Golang收入情况"})
// 	tmp := bar.AddXAxis(nameItems)
// 	for _, city := range cities {
// 		tmpArray := make([]int, 0)
// 		for _, nameItem := range nameItems {
// 			if num, exist := cityCountMap[city][nameItem]; exist {
// 				tmpArray = append(tmpArray, num)
// 			} else {
// 				tmpArray = append(tmpArray, 0)
// 			}
// 		}
// 		tmp.AddYAxis(city, tmpArray)
// 	}
// 	f, err := os.Create(path + "bar.html")
// 	if err != nil {
// 		log.Printf("show Error %s\n", err)
// 	}
// 	bar.Render(f)
// }

// BarCityJobNum return *charts.Bar
func BarCityJobNum(keywords []string) *charts.Bar {
	nameItems, count := model.BarDataCityJobNum(keywords)
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.TitleOpts{Title: "城市分布"},
		charts.ToolboxOpts{Show: true})
	bar = bar.AddXAxis(nameItems)
	for index, keyword := range keywords {
		bar.AddYAxis(keyword, count[index])
	}
	return bar
}

func BarSalaryWork(keywords []string) *charts.Bar {
	nameItems, count := model.BarDataSalaryWork(keywords)
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.TitleOpts{Title: "薪资--工作经验"},
		charts.ToolboxOpts{Show: true})
	bar = bar.AddXAxis(nameItems)
	for index, keyword := range keywords {
		bar.AddYAxis(keyword, count[index])
	}
	return bar
}

func BarSalaryEducation(keywords []string) *charts.Bar {
	nameItems, count := model.BarDataSalaryEducation(keywords)
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.TitleOpts{Title: "薪资--教育水平"},
		charts.ToolboxOpts{Show: true})
	bar = bar.AddXAxis(nameItems)
	for index, keyword := range keywords {
		bar.AddYAxis(keyword, count[index])
	}
	return bar
}

func BarSalaryCity(keywords []string) *charts.Bar {
	nameItems, count := model.BarDataSalaryCity(keywords)
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.TitleOpts{Title: "薪资--主要城市"},
		charts.ToolboxOpts{Show: true})
	bar = bar.AddXAxis(nameItems)
	for index, keyword := range keywords {
		bar.AddYAxis(keyword, count[index])
	}
	return bar
}
