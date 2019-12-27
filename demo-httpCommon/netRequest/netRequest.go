package netRequest

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
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

func PostRequest(surl string, header map[string]string, body interface{}) (result []byte) {

	client := http.Client{}
	requestBody := &strings.Reader{}

	switch body.(type) {
	case map[string]string:
		value := url.Values{}
		val := reflect.ValueOf(body).MapRange()
		for val.Next() {
			value.Add(val.Key().String(), val.Value().String())
		}
		requestBody = strings.NewReader(value.Encode())
	case string:
		requestBody = strings.NewReader(reflect.ValueOf(body).String())
	case []byte:
		requestBody = strings.NewReader(string(reflect.ValueOf(body).Bytes()))
	default:
		requestBody = strings.NewReader("")
		log.Println("Unsupported type，body is empty")
	}
	//初始 body 为 map，为了更加通用将 map 换成 interface{}
	//if body != nil {
	//	for k, v := range body{
	//		//url_v, _ := url.Parse(v)          //字符串转化为url指针
	//		value.Add(k, v)
	//		//url_v.EscapedPath() //url指针转化为字符串后加入请求体中
	//	}
	//}

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
