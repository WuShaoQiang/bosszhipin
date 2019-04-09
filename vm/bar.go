package vm

import (
	"log"
	"os"

	"github.com/chenjiandongx/go-echarts/charts"
)

var (
	path = "/home/shelljo/go/src/github.com/WuShaoQiang/crawler/boss/"
)

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