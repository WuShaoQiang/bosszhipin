package controller

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/WuShaoQiang/crawler/boss/model"

	"github.com/WuShaoQiang/crawler/boss/vm"

	"github.com/chenjiandongx/go-echarts/charts"
)

type router struct {
	name string
	charts.RouterOpts
}

var (
	basePath = "/home/shelljo/go/src/github.com/WuShaoQiang/crawler/boss"
	host     = "http://127.0.0.1:8080"
	keywords []string

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
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(basePath+"/static")))))

	http.HandleFunc("/", indexHandler)
	// http.HandleFunc("/map", mapHandler)
	http.HandleFunc("/bar", barHandler)
}

// StartUp start server
func StartUp() {
	http.ListenAndServe(addr, nil)
}

// func staticHandler() http.Handler {
// 	dir, err := os.Open("/home/shelljo/go/src/github.com/WuShaoQiang/crawler/boss/static")
// 	if err != nil {
// 		log.Fatalln("staticHandler Open File Error : ", err)
// 	}

// 	http.FileServer(dir)

// }

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/contents/index.html")
		if err != nil {
			logger.Debug("templateParseFiles Error:", err)
		}
		tmpl.Execute(w, nil)
	}
	if r.Method == http.MethodPost {
		r.ParseForm()
		keys := r.Form.Get("keywords")
		refresh := isRefresh(r.Form.Get("refresh"))
		keywords = strings.Split(keys, ",")
		log.Println("keywords : ", keywords)
		model.CrawlerGo(keywords, refresh)
		http.Redirect(w, r, "/bar", http.StatusSeeOther)
	}
}

// func mapHandler(w http.ResponseWriter, r *http.Request) {
// 	tpName := "map.html"
// 	if r.Method == http.MethodGet {
// 		templates[tpName].Execute(w, nil)
// 	}
// 	if r.Method == http.MethodPost {
// 		page := charts.NewPage(orderRouters("map")...)
// 		page.Add(
// 			vm.MapVisualMap(model.MapDataProvinceJobNum()),
// 		)
// 		page.Render(w)
// 	}

// }

func barHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		page := charts.NewPage(orderRouters("bar")...)
		page.Add(
			vm.BarCityJobNum(keywords),
		)
		page.Render(w)
	}
}
