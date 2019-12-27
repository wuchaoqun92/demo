package main

import (
	"demo-person/demo-OpenApiTest/imageRecognitionBaidu/imageEncode"
	"demo-person/demo-httpCommon/netRequest"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const access_token = "you access_token"

var car_url = "https://aip.baidubce.com/rest/2.0/image-classify/v1/car?access_token="

type car_info struct {
	Name       string  `json:"name"`
	Year       string  `json:"year"`
	Score      float64 `json:"score"`
	Baike_info baike_info
}

type baike_info struct {
	Baike_url   string `json:"baike_url"`
	Image_url   string `json:"image_url"`
	Description string `json:"description"`
}

type Res struct {
	Result []car_info `json:"result"`
}

func main() {

	token := access_token
	requestUrl := car_url + token

	image_byte, _ := imageEncode.LoadImage("/Users/wuchaoqun/Desktop/1.jpg")
	image_value := base64.StdEncoding.EncodeToString(image_byte)

	header := make(map[string]string, 0)
	header["Content-Type"] = "application/x-www-form-urlencoded"

	body := make(map[string]string, 0)
	body["image"] = string(image_value)
	body["top_num"] = "5"
	body["baike_num"] = "1"

	result := netRequest.PostRequest(requestUrl, header, body)

	fmt.Println("result is", string(result))

	res := Res{}

	err := json.Unmarshal(result, &res)
	if err != nil {
		return
	}
	fmt.Println(res)

	if res.Result[0].Score > 0.99999 {
		resp, err := http.Get(res.Result[0].Baike_info.Image_url)
		if err != nil {
			fmt.Println("get picture failed,err is", err)
			return
		}
		by, _ := ioutil.ReadAll(resp.Body)

		file, err := os.Create("/Users/wuchaoqun/Desktop/x.jpg")
		n, _ := file.Write(by)
		fmt.Println("sucess,write len:", n)
	} else {
		fmt.Println("failed,no car  ")
	}

}
