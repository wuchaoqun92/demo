package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func main() {

	test2()

}
func test2() {
	resp, err := http.Get("https://mp.weixin.qq.com/s/I2XbTrYkfjGEfqSnhHFB3g")
	if err != nil {
		fmt.Println("get picture failed,err is", err)
		return
	}
	by, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	commentCount := `<div id="img-content"><h2.*?</div> ` //`<h2 class="rich_media_title" id="activity-name">([^<]+)</h2>`
	reg := regexp.MustCompile(commentCount)

	//fmt.Println(string(by))

	for i, d := range reg.FindAllString(string(by), -1) {
		fmt.Println(i, d)
	}

}

func test1() {
	resp, err := http.Get("https://mp.weixin.qq.com/s/F6Hao3sL8UGDmRI6BXskCQ")
	if err != nil {
		fmt.Println("get picture failed,err is", err)
		return
	}
	by, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	commentCount := `<span (.*?)>(.*?)</span>`
	reg := regexp.MustCompile(commentCount)

	res := make([]string, 0)
	for _, d := range reg.FindAllString(string(by), -1) {
		n := Utf8Index(d, ">")
		res = append(res, d[n:])

	}
	res1 := make([]string, len(res))
	commentCount1 := `data-src="(.*?)"`
	reg1 := regexp.MustCompile(commentCount1)
	for i, v := range res {
		//fmt.Println(v)
		a := strings.ReplaceAll(v, "><strong>", "小标题：")
		d := strings.ReplaceAll(a, "</strong></span(.*?)>", "")
		e := strings.ReplaceAll(d, "</span>", "")
		//b := strings.ReplaceAll(e,">","正文：")
		res1[i] = e
		c := reg1.FindString(e)
		if c != "" {
			res1[i] = c
		}

	}
	for _, v := range res1 {
		fmt.Println(v)
		//strings.ReplaceAll(v,">","正文")
	}
}

//获取字符在字符串中的位置
func Utf8Index(str, substr string) int {
	asciiPos := strings.Index(str, substr)
	if asciiPos == -1 || asciiPos == 0 {
		return asciiPos
	}
	pos := 0
	totalSize := 0
	reader := strings.NewReader(str)
	for _, size, err := reader.ReadRune(); err == nil; _, size, err = reader.ReadRune() {
		totalSize += size
		pos++
		// 匹配到
		if totalSize == asciiPos {
			return pos
		}
	}
	return pos
}
