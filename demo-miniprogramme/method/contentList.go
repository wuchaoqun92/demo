package method

import (
	"awesomeProject/demo-miniprogramme/common"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ContentList struct {
	Code string `json:"code" db:"code"`
	ListTitle string	`json:"list_title" db:"list_title"`
	ContentId int `json:"content_id" db:"list_id"`
	Counts int 	`json:"counts" db:"counts"`
}

type Content struct {
	ContentList
	Title string `json:"title" db:"title"`
	Author int `json:"author" db:"user_id"`
	CreateTime string `json:"create_time" db:"create_time"`
	Text string `json:"text" db:"text"`
	ContentTime string `json:"content_time" db:"content_time"`
}

func ChoiceMethodToList(w http.ResponseWriter,r *http.Request){

	var RcvMsg common.RecMessage
	var Backmsg common.BackMessage
	var list []Content
	var err error

	str := r.FormValue("msg")
	fmt.Println("recv data is",str)

	RcvMsg.ReadMsg(str)
	fmt.Println(RcvMsg)

	switch RcvMsg.Cmd {
	case common.GetConList:
		fmt.Println("获取全部列表")
		list,err = GetConlist(RcvMsg.Data)
	case common.GetConDetail:
		fmt.Println("获取内容详情")
		list,err=GetConDetail(RcvMsg.Data)
	case common.InsertConDetail:
		fmt.Println("插入内容详情")
		err=InsertConDetail(RcvMsg.Data)
	default:
		fmt.Print("world")

	}

	if err != nil {
		fmt.Println("获取数据错误，err：",err)
		dataBack,_ := Backmsg.WriteMsg("0","获取数据错误")
		io.WriteString(w,dataBack)
		return
	}

	da,err := json.Marshal(list)
	fmt.Println("back data is",string(da))
	dataBack,_ := Backmsg.WriteMsg("1","获取数据成功",string(da))
	io.WriteString(w,dataBack)

}


func GetConlist(code string)(content_lst []Content,err error){


	list := GetListFromDatabase(code)

	content_lst = make([]Content,len(list))

	for i,v := range list{
		content_lst[i].Code = v.Code
		content_lst[i].ListTitle = v.ListTitle
		content_lst[i].ContentId = v.ContentId
		content_lst[i].Counts = v.Counts
	}

		if len(content_lst)<0 {
			err = errors.New("无数据")
			return
		}

	return
	}

func GetConDetail(content_id string)(content_detail []Content,err error){

	list := GetContentDetailFromDatabase(content_id)
	if len(list)<0 {
		err = errors.New("无数据")
		return
	}

	content_detail = make([]Content,len(list))

	for i,v := range list{
		content_detail[i].Title = v.Title
		content_detail[i].CreateTime = v.CreateTime
		content_detail[i].Text = v.Text
		content_detail[i].ContentTime = v.ContentTime
	}

	return
}

func InsertConDetail(MsgData string)(err error){

	var contentDatail = Content{}
	err = json.Unmarshal([]byte(MsgData),&contentDatail)
	if err != nil {
		fmt.Println("receive to insert data unmarshal failed",err)
		return
	}


	return
}
