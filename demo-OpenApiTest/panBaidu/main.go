package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	ACCESS_TOKEN = "your access_token"
)

type PicList struct {
	Fs_id           uint64 `json:"fs_id"`           //文件在云端的唯一标识ID
	Path            string `json:"path"`            //文件的绝对路径
	Server_filename string `json:"server_filename"` //文件名称
	Size            int    `json:"size"`            //文件大小，单位B
}

type List struct {
	Errno      int       `json:"errno"`
	Guid_info  string    `json:"guid_info"`
	Info       []PicList `json:"list"`
	Request_id int       `json:"request_id"`
	Guid       int       `json:"guid"`
}

func main() {
	//url := "https://pan.baidu.com/rest/2.0/xpan/multimedia?method=filemetas&access_token="+ACCESS_TOKEN+"&fs_id=1071043601441656"
	url := "https://pan.baidu.com/rest/2.0/xpan/file?method=list&access_token=" + ACCESS_TOKEN + "&page=1&num=10&dir=/图片pic"

	client := &http.Client{}
	//生成要访问的url
	//url := "http://somesite/somepath/"

	//提交请求
	reqest, err := http.NewRequest("GET", url, nil)

	//增加header选项
	//reqest.Header.Add("Cookie", "xxxxxx")
	reqest.Header.Add("User-Agent", "pan.baidu.com")
	//reqest.Header.Add("X-Requested-With", "xxxx")

	if err != nil {
		panic(err)
	}
	//处理返回结果
	resp, _ := client.Do(reqest)
	defer resp.Body.Close()

	if err != nil {
		fmt.Println("http.Get err=", err)
		return
	}

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("ioutil.ReadAll err=", err)
		return
	}
	fmt.Println("bytes is", string(bytes))

	list := List{}

	err = json.Unmarshal(bytes, &list)
	if err != nil {
		fmt.Println("unmarshal faile:", err)
		return
	}

	fmt.Println("list is", list)
	for _, v := range list.Info {
		fmt.Printf("文件名:%s,文件id:%d\n", v.Server_filename, v.Fs_id)
	}

}

func unicodeToChinese(sText string) {
	textQuoted := strconv.QuoteToASCII(sText)
	textUnquoted := textQuoted[1 : len(textQuoted)-1]
	fmt.Println(textUnquoted)

	sUnicodev := strings.Split(textUnquoted, "\\u")
	var context string
	for _, v := range sUnicodev {
		if len(v) < 1 {
			continue
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			panic(err)
		}
		context += fmt.Sprintf("%c", temp)
	}
	fmt.Println(context)
}

type ab int

func (t ab) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("aaaaa")
}

func a() {

	var ab ab
	var handle http.Handler
	handle = ab

	s := &http.Server{
		Addr:           ":8080",
		Handler:        handle,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
