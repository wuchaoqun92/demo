package token

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	Api_key     = "your api_key"
	Secret_key  = "your secret_key"
	Request_url = "https://aip.baidubce.com/oauth/2.0/token"
)

type Token struct {
	Access_token string    `json:"access_token"`
	Expires_in   int       `json:"expires_in"`
	ValidityTime time.Time //有效期截止时间
	ErrNo        bool      //判断token是否正确取得
}

func GetToken() (token Token, err error) {

	url := Request_url + "?grant_type=client_credentials&client_id=" + Api_key + "&client_secret=" + Secret_key

	client := http.Client{}

	reqest, err := http.NewRequest("POST", url, nil)

	reqest.Header.Add("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		panic(err)
	}
	//处理返回结果
	resp, _ := client.Do(reqest)
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(res, &token)
	if err != nil {
		fmt.Println("ummarshal failed,err is", err)
		token.ErrNo = false
		return
	}

	time := time.Now().Add(time.Duration(token.Expires_in) * time.Second)
	token.ValidityTime = time
	token.ErrNo = true

	return
}
