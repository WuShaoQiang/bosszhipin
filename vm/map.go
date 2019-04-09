package vm

import "github.com/chenjiandongx/go-echarts/charts"

func MapVisualMap(mapData map[string]float32) *charts.Map {
	mc := charts.NewMap("china")
	mc.SetGlobalOptions(
		charts.TitleOpts{Title: "Map-设置 VisualMap"},
		charts.VisualMapOpts{Calculable: true},
	)
	mc.Add("map", mapData)
	return mc
}
