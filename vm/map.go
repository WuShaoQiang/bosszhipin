package vm

import (
	"fmt"

	"github.com/WuShaoQiang/crawler/boss/model"
	"github.com/chenjiandongx/go-echarts/charts"
)

//MapVisualMap return *charts.Map
func MapVisualMap(keyword string) *charts.Map {
	mapData := make(map[string]float32)
	mc := charts.NewMap("china")
	mc.SetGlobalOptions(
		charts.TitleOpts{Title: fmt.Sprintf("Map - %s", keyword)},
		charts.VisualMapOpts{Calculable: true},
	)
	mapData = model.MapDataProvinceJobNum(keyword)
	mc.Add("map", mapData)
	return mc
}
