package controller

import (
	"log"
	"net/http"

	"github.com/WuShaoQiang/crawler/boss/model"

	"github.com/WuShaoQiang/crawler/boss/vm"

	"github.com/chenjiandongx/go-echarts/charts"
)

type router struct {
	name string
	charts.RouterOpts
}

var (
	path = "/home/shelljo/go/src/github.com/WuShaoQiang/crawler/boss/"
	host = "http://127.0.0.1:8080"

	routers = []router{
		{"bar", charts.RouterOpts{URL: host + "/bar", Text: "Bar-(柱状图)"}},
		{"bar3D", charts.RouterOpts{URL: host + "/bar3D", Text: "Bar3D-(3D 柱状图)"}},
		{"boxPlot", charts.RouterOpts{URL: host + "/boxPlot", Text: "BoxPlot-(箱线图)"}},
		{"effectScatter", charts.RouterOpts{URL: host + "/effectScatter", Text: "EffectScatter-(动态散点图)"}},
		{"funnel", charts.RouterOpts{URL: host + "/funnel", Text: "Funnel-(漏斗图)"}},
		{"gauge", charts.RouterOpts{URL: host + "/gauge", Text: "Gauge-仪表盘"}},
		{"geo", charts.RouterOpts{URL: host + "/geo", Text: "Geo-地理坐标系"}},
		{"graph", charts.RouterOpts{URL: host + "/graph", Text: "Graph-关系图"}},
		{"heatMap", charts.RouterOpts{URL: host + "/heatMap", Text: "HeatMap-热力图"}},
		{"kline", charts.RouterOpts{URL: host + "/kline", Text: "Kline-K 线图"}},
		{"line", charts.RouterOpts{URL: host + "/line", Text: "Line-(折线图)"}},
		{"line3D", charts.RouterOpts{URL: host + "/line3D", Text: "Line3D-(3D 折线图)"}},
		{"liquid", charts.RouterOpts{URL: host + "/liquid", Text: "Liquid-(水球图)"}},
		{"map", charts.RouterOpts{URL: host + "/map", Text: "Map-(地图)"}},
		{"overlap", charts.RouterOpts{URL: host + "/overlap", Text: "Overlap-(重叠图)"}},
		{"parallel", charts.RouterOpts{URL: host + "/parallel", Text: "Parallel-(平行坐标系)"}},
		{"pie", charts.RouterOpts{URL: host + "/pie", Text: "Pie-(饼图)"}},
		{"radar", charts.RouterOpts{URL: host + "/radar", Text: "Radar-(雷达图)"}},
		{"sankey", charts.RouterOpts{URL: host + "/sankey", Text: "Sankey-(桑基图)"}},
		{"scatter", charts.RouterOpts{URL: host + "/scatter", Text: "Scatter-(散点图)"}},
		{"scatter3D", charts.RouterOpts{URL: host + "/scatter3D", Text: "Scatter-(3D 散点图)"}},
		{"surface3D", charts.RouterOpts{URL: host + "/surface3D", Text: "Surface3D-(3D 曲面图)"}},
		{"themeRiver", charts.RouterOpts{URL: host + "/themeRiver", Text: "ThemeRiver-(主题河流图)"}},
		{"wordCloud", charts.RouterOpts{URL: host + "/wordCloud", Text: "WordCloud-(词云图)"}},
		{"page", charts.RouterOpts{URL: host + "/page", Text: "Page-(顺序多图)"}},
	}
)

// Register register all handlers
func Register() {
	http.HandleFunc("/map", mapHandler)
	http.HandleFunc("/bar", barHandler)
}

func mapHandler(w http.ResponseWriter, _ *http.Request) {
	page := charts.NewPage(orderRouters("map")...)
	page.Add(
		vm.MapVisualMap(model.MapDataProvinceJobNum()),
	)
	page.Render(w)
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	tpName := "bar.html"
	if r.Method == http.MethodGet {
		templates[tpName].Execute(w, nil)
	}
	if r.Method == http.MethodPost {
		var keywords []string
		r.ParseForm()
		keyword1 := r.Form.Get("keyword1")
		keyword2 := r.Form.Get("keyword2")
		keywords = append(keywords, keyword1, keyword2)
		log.Println("User post keywords", keywords)

		page := charts.NewPage(orderRouters("bar")...)
		page.Add(
			vm.BarCityJobNum(keywords),
		)
		page.Render(w)
	}

}
