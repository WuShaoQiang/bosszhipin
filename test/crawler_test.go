package test

import (
	"net/http"
	"net/url"
	"testing"
)

func TestCrawler(t *testing.T) {
	// keywords := []string{"golang", "java"}
	keywords := "golang,java"
	refresh := "on"
	value := url.Values{
		keywords: []string{keywords},
		refresh:  []string{refresh},
	}
	resp, err := http.PostForm("http://127.0.0.1:8080", value)
	if err != nil {
		t.Fatal("PostForm Error : ", err)
	}

	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != 200 {
		t.Errorf("Receive a statuscode : %v", resp.StatusCode)
	}
}
