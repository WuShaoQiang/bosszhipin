package model

import (
	"math/rand"
)

var uas = [...]string{
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/38.0.2125.111 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1",
	"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/536.3",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/535.24",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_0) AppleWebKit/536.3",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko)",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11",
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.103 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; rv:11.0) like Gecko",
}

// Proxy stores the proxy filtered from CrudeProxy
type Proxy struct {
	// ID is the ID value of the current record, which is unique among all proxies.
	ID int64 `gorm:"AUTO_INCREMENT;" json:"id"`
	// IP is the IP address of the proxy. e.g 127.0.0.1
	IP string `json:"ip"`
	// Port is the Port of the proxy. e.g 3306
	Port string `json:"port"`
	// SchemeType represents the protocol type supported by the proxy.
	// 0: http
	// 1: https
	// 2: http & https
	SchemeType int64 `json:"scheme_type"`
	// Content is the ip:port of the proxy. e.g 127.0.0.1:3306
	Content string `gorm:"unique_index:unique_content;" json:"content"`

	// AssessTimes is the number of evaluations of the proxy
	AssessTimes int64 `json:"assess_times"`
	// SuccessTimes is the number of successful evaluations of the proxy
	SuccessTimes int64 `json:"success_times"`
	// AvgResponseTime is the average response time of the proxy
	AvgResponseTime float64 `json:"avg_response_time"`
	// ContinuousFailedTimes is the number of consecutive failures during the proxy evaluation process
	ContinuousFailedTimes int64 `json:"continuous_failed_times"`
	// Score is the rating of the proxy
	Score float64 `json:"score"`
	// InsertTime is the insertion time of the proxy
	InsertTime int64 `json:"insert_time"`
	// UpdateTime is the update time of the proxy, can also reflect the last evaluation time
	UpdateTime int64 `json:"update_time"`
}

//GetIP return IP:Port randomly
func GetIP() (addr string) {
	// addr = ""
	rows, err := db.Table("proxy").Where("score > ?", 5).Select("content").Rows()
	if err != nil {
		logger.Debugln("GetIP Error : ", err)
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&addr)
	}
	return
}

// DeleteIP delete ip in mysql
func DeleteIP(addr string) {
	err := db.Table("proxy").Where("content=?", addr).Delete(&Proxy{}).Error
	if err != nil {
		logger.Debugln("DeleteIP Error : ", err)
	}
}

// GetUserAgent return
func GetUserAgent() string {
	n := rand.Intn(len(uas))
	return uas[n]
}
