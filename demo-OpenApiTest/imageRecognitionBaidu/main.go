package main

import (
	"demo-person/demo-OpenApiTest/imageRecognitionBaidu/imageEncode"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	//a,b:=token.GetToken()
	//c := a.ValidityTime.Sub(time.Now())
	//if b == nil{
	//	if c > 0 {
	//		fmt.Println(a.Access_token,a.ValidityTime,c)
	//	}else {
	//		fmt.Println(c)
	//	}
	//}
	token := "you access_token"
	bytes := imageEncode.LoadImage("/Users/wuchaoqun/Desktop/1.jpg")
	_, image_value := imageEncode.ImageEncodeToBase64(bytes)
	getResult(token, image_value)
	//fmt.Println(token,base64code)

}

func getResult(accessToken string, image *url.URL) {
	//image,_ = urlEncodeAndDecode(image)
	//fmt.Println("urlencode:",image)
	requestUrl := "https://aip.baidubce.com/rest/2.0/image-classify/v1/car?access_token=" + accessToken
	//surl :=requestUrl+"?access_token="+accessToken+"&image="+image+"&top_num="+"5"

	value := url.Values{}
	value.Add("image", image.EscapedPath())
	value.Add("top_num", "5")

	//resp,_:=http.PostForm(requestUrl,value)

	body := strings.NewReader(value.Encode())

	client := http.Client{}

	reqest, err := http.NewRequest("POST", requestUrl, body)
	//
	reqest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//reqest.PostForm = url.Values{}
	//reqest.PostForm.Add("image",image.EscapedPath())
	//reqest.PostForm.Add("top_num","5")

	//处理返回结果
	//resp,err := client.PostForm(requestUrl,value)
	if err != nil {
		panic(err)
	}
	resp, _ := client.Do(reqest)

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("ioutil.ReadAll err=", err)
		return
	}

	fmt.Println("result is", string(bytes))
}

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
