package method

import (
	"demo-person/demo-miniprogramme/common"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CheckIn(w http.ResponseWriter, r *http.Request) {

	var checkCode common.CheckIn
	var msg common.BackMessage
	var dataBack string

	str := r.FormValue("season")
	fmt.Println("recv data is", str)

	err := json.Unmarshal([]byte(str), &checkCode)
	if err != nil {
		fmt.Println("unmarshale err:", err)
	}

	openid := GetOpenId(checkCode.Code)
	if openid.Cmd == "0" {
		fmt.Println("获取openid失败")

		dataBack, err := msg.WriteMsg("0", "获取openid失败")
		if err != nil {
			return
		}
		fmt.Println("back msg is", msg)
		io.WriteString(w, dataBack)
		return
	}

	if CheckCode(checkCode.Num, openid.Data) != nil {
		dataBack, err = msg.WriteMsg("0", "验证失败", openid.Data)
	} else {
		dataBack, err = msg.WriteMsg("1", "成功")
	}

	fmt.Println("back msg is", msg)

	io.WriteString(w, dataBack)
}
