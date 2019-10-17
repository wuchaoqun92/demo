package common

import (
	"encoding/json"
	"fmt"
)

type CheckIn struct {
	Num string 	`json:"num"`
	Code string `json:"code"`
}

type RecMessage struct {
	Cmd string `json:"cmd"`
	Data string `json:"data"`
}

type BackMessage struct {
	Cmd string `json:"cmd"` //类型状态
	Data string `json:"data"` //返回数据
	Msg string `json:"msg"` //提示信息
}


func (msg *RecMessage)ReadMsg(str string)(err error){

	err = json.Unmarshal([]byte(str),msg)
	if err != nil {
		fmt.Println("RecMessage unmarshale err:",err)
	}

	return
}

func (msg *BackMessage)WriteMsg(args... string)(databack string,err error){

	msg.Cmd = args[0]
	msg.Msg = args[1]

	for i:=2;i<len(args) ;i++  {
		msg.Data += args[i]
	}

	da,err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Back data marshal failed,err is",err)
	}
	databack = string(da)

	return

}

