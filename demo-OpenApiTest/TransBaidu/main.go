package main

import (
	"crypto/md5"
	myclip "demo-person/demo-OpenApiTest/TransBaidu/clipboard"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	appid = "your appid"
	key   = "your key"
)

type result struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}
type back struct {
	From         string   `json:"from"`
	To           string   `json:"to"`
	Trans_result []result `json:"trans_result"`
}

type language struct {
	Error_code interface{} `json:"error_code"`
	Error_msg  string      `json:"error_msg"`
	Data       langData
}
type langData struct {
	Src string `json:"src"`
}

var lasttime string
var languages string

func init() {
	lasttime = getClipboard()
}

func main() {

	//query := strings.ReplaceAll(getClipboard()," ","+")
	//fmt.Printf("result is:'%s'\n",getResult(query))

	c := cron.New()
	spec := "*/5 * * * * *"
	c.AddFunc(spec, func() {
		//fmt.Println("unicode is",strconv.QuoteToASCII(getClipboard()))

		if lasttime != getClipboard() {
			q := checkQuery(getClipboard())
			query := strings.Split(q, "\n")
			fmt.Println("----RESULT IS:")
			to, from := languageCheckRes(q)
			if len(query) < 5 {
				for _, v := range query {
					fmt.Printf("译文：%s\n", getResult(checkQuery(v), to, from))
					time.Sleep(time.Second)
				}
			} else {
				fmt.Println("内容超过5行，暂时无法翻译。11-1之后可用")
			}

			//fmt.Printf("result is:'%s'\n",getResult(getClipboard()))
			lasttime = getClipboard()
		}
	})

	go c.Start()
	defer c.Stop()

	select {
	//case <-time.After(time.Second * 10):
	//	return
	}

}

//语种识别QPS建议控制在3次以下，语句翻译QPS以升级至10。故语种判断针对整段内容进行
func languageCheckRes(q string) (to, from string) {
	to = "zh"
	from = "auto"
	//q = checkQuery(q)

	switch languageCheck(q) {
	case "zh":
		from = "zh"
		to = "en"
		fmt.Println("原中文:", q)
	case "en":
		from = "en"
		to = "zh"
		fmt.Println("Original English:", q)
	default:
		fmt.Println("其他语种:", q)
	}
	return
}

//制作url，并获取返回值。同时根据返回值进行反序列化获得需要的结果
func getResult(q, to, from string) (result string) {

	url := makeUrl(q, to, from)
	res := Req(url)

	result1 := back{}
	err := json.Unmarshal(res, &result1)
	if err != nil {
		fmt.Println("ummaishal failed,err is", err)
		return
	}
	//fmt.Println(string(res))
	for _, v := range result1.Trans_result {
		result += v.Dst
	}
	return //result1.Trans_result[0].Dst
}

//字符串处理。部分符号无法识别
func checkQuery(q string) string {
	q = strings.TrimSpace(q)
	q = strings.TrimPrefix(q, "A.")
	q = strings.TrimPrefix(q, "B.")
	q = strings.TrimPrefix(q, "C.")
	q = strings.TrimPrefix(q, "D.")
	q = strings.ReplaceAll(q, "-", ";")

	//query := strings.Split(q,"\n")
	//for i,v := range query  {
	//	fmt.Printf("key : %d,value : %s\n",i,v)
	//}
	return q
}

//生成待访问的url
func makeUrl(q, to, from string) (url string) {
	from = from
	to = to
	appid := appid
	salt := fmt.Sprintf("%d", rand.Intn(1000000000))
	key := key
	sign := md5V(appid + q + salt + key)

	query, _ := urlEncodeAndDecode(q)
	//fmt.Println("q urlencode 前",q)
	//fmt.Println("q urlencode 后",query)

	url = "http://api.fanyi.baidu.com/api/trans/vip/translate?q=" + query + "&from=" + from + "&to=" + to + "&appid=" + appid + "&salt=" + salt + "&sign=" + sign
	//fmt.Println("请求的url：",url)
	return
}

//带翻译内容的语种辨别
func languageCheck(q string) (res string) {
	appid := appid
	salt := fmt.Sprintf("%s", rand.Intn(1000000000))
	key := key
	sign := md5V(appid + q + salt + key)

	query, _ := urlEncodeAndDecode(q)

	url := "https://fanyi-api.baidu.com/api/trans/vip/language?q=" + query + "&appid=" + appid + "&salt=" + salt + "&sign=" + sign
	resp, _ := http.Get(url)

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("ioutil.ReadAll err=", err)
		return
	}
	//fmt.Println("language data is",string(bytes))

	languages := language{}
	err = json.Unmarshal(bytes, &languages)
	if err != nil {
		fmt.Println("ummarshal failed,err is", err)
		return
	}
	//fmt.Println(string(res))
	res = languages.Data.Src

	return

}

//访问url，获取返回值
func Req(url string) (bytes []byte) {

	client := http.Client{}

	reqest, err := http.NewRequest("POST", url, nil)

	reqest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		panic(err)
	}
	//处理返回结果
	resp, _ := client.Do(reqest)
	defer resp.Body.Close()

	bytes, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("ioutil.ReadAll err=", err)
		return
	}
	//fmt.Println("bytes is",string(bytes))
	return
}

func getClipboard() (content string) {
	//clipboard.WriteAll(`复制这段内容到剪切板`)
	// 读取剪切板中的内容到字符串
	content, err := myclip.ReadAll()
	if err != nil {
		panic(err)
	}
	//println(content)
	return
}

//MD5序列化字符串
func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

//url encode/decode，url字符类型编译与反编译
func urlEncodeAndDecode(a string) (encode string, decode url.Values) {
	v := url.Values{}
	v.Add("key", a)
	encode = v.Encode()
	encode = strings.TrimPrefix(encode, "key=")
	//fmt.Println(v)
	//fmt.Println(encode)
	// url decode
	decode, _ = url.ParseQuery(encode)
	//fmt.Println(decode)
	return
}
