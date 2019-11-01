package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	resp, err := http.Get("https://mp.weixin.qq.com/s/F6Hao3sL8UGDmRI6BXskCQ")
	if err != nil {
		fmt.Println("get picture failed,err is", err)
		return
	}
	by, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	//fmt.Println(string(by))

	file, err := os.Create("/Users/wuchaoqun/Desktop/x.html")
	defer file.Close()
	//n,_ :=file.Write(by)
	//fmt.Println("sucess,write len:",n)

	commentCount := `<span style="(.*?)">(.*?)</span>`
	//commentCount := `meta name="(.*?)" content="(.*?)"`
	reg := regexp.MustCompile(commentCount)
	//txt2 := rp2.FindAllString(string(by), -1)
	//fmt.Println(txt2)
	res := make([]string, 0)
	for _, d := range reg.FindAllString(string(by), -1) {
		//fmt.Println(i,d)
		n := Utf8Index(d, ">")
		res = append(res, d[n:])

	}
	res1 := make([]string, len(res))
	commentCount1 := `data-src="(.*?)"`
	reg1 := regexp.MustCompile(commentCount1)
	for i, v := range res {
		//fmt.Println(v)
		a := strings.ReplaceAll(v, "><strong>", "小标题：")
		d := strings.ReplaceAll(a, "</strong></span>", "")
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
