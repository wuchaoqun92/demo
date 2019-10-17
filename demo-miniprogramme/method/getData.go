package method

import (
	"encoding/json"
	"errors"
	"fmt"
)

type GetData interface {
	GetData()(data []byte,err error)
}


func (p *ContentList)GetData()(data []byte,err error){

	list := GetListFromDatabase(p.Code)

	content_lst := make([]Content,len(list))

	for i,v := range list{
		content_lst[i].Code = v.Code
		content_lst[i].ListTitle = v.ListTitle
		content_lst[i].ContentId = v.ContentId
	}
	if len(content_lst)<0 {
		err = errors.New("无数据")
		return
	}
	data,err = json.Marshal(content_lst)
	return

}

func a(){
	var p *ContentList
	p.Code="3"
	var x GetData

	x=p
	x.GetData()
	fmt.Println(x)
}