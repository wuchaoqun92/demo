package method

import (
	"demo-person/demo-miniprogramme/common"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	appid      = "wx66eda95cb0ef5ae4"
	secret     = "267ecc70e7aee3098a47a52798d8c628"
	grant_type = "authorization_code"
)

type OId struct {
	Session_key string `json:"session_key"`
	Openid      string `json:"openid"`
}

type errMessage struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

//GET https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code

func GetOpenId(code string) (msg common.BackMessage) {

	url := "https://api.weixin.qq.com/sns/jscode2session?appid=" + appid + "&secret=" + secret + "&js_code=" + code + "&grant_type=" + grant_type + ""

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("http.Get err=", err)
		return
	}

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("ioutil.ReadAll err=", err)
		return
	}
	fmt.Println(string(bytes))

	var oid OId = OId{}
	var error errMessage = errMessage{}

	err = json.Unmarshal(bytes, &oid)
	if err != nil {
		fmt.Println("unmarshal failed,err:", err)
		return
	}

	err = json.Unmarshal(bytes, &error)
	if err != nil {
		fmt.Println("unmarshal failed,err:", err)
		return
	}

	if oid.Openid == "" {
		msg = common.BackMessage{
			Cmd:  "0",
			Msg:  error.Errmsg,
			Data: string(error.Errcode),
		}
		return
	}
	msg = common.BackMessage{
		Cmd:  "1",
		Msg:  "成功",
		Data: oid.Openid,
	}
	return
}
