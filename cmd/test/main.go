package main

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	keywords = []string{"golang实习", "golang", "后端开发实习生"}
)

func main() {
	for _, keyword := range keywords {
		// chinese := getChinese(keyword)
		// fmt.Println(chinese)
		// encoded := urlEncode(chinese)
		// fmt.Println(encoded)
		fmt.Println(urlEncode(keyword))
	}
}

func urlEncode(keyword string) string {
	reg := regexp.MustCompile(`[\p{Han}]+`)
	strs := reg.FindAllString(keyword, -1)
	chinese := strings.Join(strs, "")
	encoded := chineseEncode(chinese)
	encodedURL := reg.ReplaceAllString(keyword, encoded)
	return encodedURL
}

func chineseEncode(chinese string) (encoded string) {
	encoded = ""
	byteChinese := []byte(chinese)
	// return fmt.Sprintf("%x")
	for _, singleByte := range byteChinese {
		encoded = encoded + "%" + fmt.Sprintf("%x", singleByte)
	}
	return
}
