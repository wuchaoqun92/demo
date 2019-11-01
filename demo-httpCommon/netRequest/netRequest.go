package netRequest

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func GetRequest(url string) (result []byte) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("getRequest failed,err is", err)
		return
	}

	result, err = ioutil.ReadAll(resp.Body)
	return
}

func PostRequest(surl string, header, body map[string]string) (result []byte) {

	client := http.Client{}

	value := url.Values{}
	if body != nil {
		for k, v := range body {
			url_v, _ := url.Parse(v)          //字符串转化为url指针
			value.Add(k, url_v.EscapedPath()) //url指针转化为字符串后加入请求体中
		}
	}

	requestBody := strings.NewReader(value.Encode())

	request, err := http.NewRequest("POST", surl, requestBody)
	if err != nil {
		log.Fatal("newRequest failed ,err is", err)
		return
	}
	if header != nil {
		for k, v := range header {
			request.Header.Add(k, v)
		}
	}

	resp, err := client.Do(request)
	defer resp.Body.Close()

	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("reader all failed,err is", err)
		return
	}
	return
}
