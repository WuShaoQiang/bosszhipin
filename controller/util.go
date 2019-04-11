package controller

import (
	"github.com/chenjiandongx/go-echarts/charts"
)

func orderRouters(chartType string) []charts.RouterOpts {
	for i := 0; i < len(routers); i++ {
		if routers[i].name == chartType {
			routers[i], routers[0] = routers[0], routers[i]
			break
		}
	}

	rs := make([]charts.RouterOpts, 0)
	for i := 0; i < len(routers); i++ {
		rs = append(rs, routers[i].RouterOpts)
	}
	return rs
}

func isRefresh(refresh string) bool {
	if refresh == "on" {
		return true
	}
	return false
}

// func populateTemplates() map[string]*template.Template {
// 	const basePath = "templates"
// 	result := make(map[string]*template.Template)

// 	layout := template.Must(template.ParseFiles(basePath + "/_base.html"))
// 	dir, err := os.Open(basePath + "/content")
// 	if err != nil {
// 		panic("Failed to open template blocks directory: " + err.Error())
// 	}

// 	fis, err := dir.Readdir(-1)
// 	if err != nil {
// 		panic("Failed to read contents of content directory: " + err.Error())
// 	}

// 	for _, fi := range fis {
// 		f, err := os.Open(basePath + "/content/" + fi.Name())
// 		if err != nil {
// 			panic("Failed to open template " + fi.Name())
// 		}
// 		content, err := ioutil.ReadAll(f)
// 		if err != nil {
// 			panic("Failed to read content from file " + fi.Name())
// 		}
// 		f.Close()
// 		tmpl := template.Must(layout.Clone())
// 		_, err = tmpl.Parse(string(content))
// 		if err != nil {
// 			panic("Failed to parse contents of" + fi.Name())
// 		}
// 		result[fi.Name()] = tmpl
// 	}
// 	return result
// }
