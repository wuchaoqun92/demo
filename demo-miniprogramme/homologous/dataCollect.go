package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"sync"
)

func DatabaseConnect()(db *sqlx.DB,err error) {
	database,err := sqlx.Open("mysql","root:12345678@tcp(127.0.0.1:3306)/miniProgramme")
	if err != nil {
		fmt.Println("database open failed,err:",err)
		return
	}

	db = database
	return
}


type dataCollect interface {
	CollectMsg(string)([]byte)
}


type titleList struct {
	Id int `db:"list_id" json:"id"`
	Title string `db:"list_title" json:"title"`
	Code string `db:"code" json:"code"`
	Counts int `db:"counts" json:"counts"`
}
type titleListMsg struct {
	titleListMsg []titleList
}

type contentList struct {
	ListId int `db:"list_id" json:"list_id"`
	Title string `db:"title" json:"title"`
	Text string `db:"text" json:"text"`
	ContentTime string `db:"content_time" json:"content_time"`
}
type contentListMsg struct {
	contentListMsg []contentList
}


func (list *titleListMsg)CollectMsg(code string)(result []byte){

	db,err := DatabaseConnect()
	if err != nil {
		fmt.Println("database connect failed,err:",err)
		return
	}
	defer db.Close()

	str :=`SELECT code,list_title,list_id,COUNT(content_id) as counts FROM
			(
			SELECT a.code,a.list_id,a.list_title,b.content_id
			FROM list a LEFT JOIN content b ON a.list_id = b.list_id
			WHERE a.code = ? 
					) c 
			GROUP BY list_id,code,list_title
			ORDER BY list_id`

	err = db.Select(&list.titleListMsg,str,code)
	if err != nil {
		fmt.Println("数据插入结构体失败，err",err)
		return
	}

	fmt.Println("gotdata is",list.titleListMsg)
	result,err = json.Marshal(list.titleListMsg)
	if err != nil {
		fmt.Println("marshal failed,err is",err)
		return
	}
	return
}


func (list *contentListMsg)CollectMsg(listId string)(result []byte){

	db,err := DatabaseConnect()
	if err != nil {
		fmt.Println("database connect failed,err:",err)
		return
	}
	defer db.Close()


	err = db.Select(&list.contentListMsg, "SELECT list_id,text,title,content_time FROM content where list_id = ? ",listId)
	if err != nil {
		fmt.Println("数据插入结构体失败，err:",err)
		return
	}
	fmt.Println("gotdata is",list)
	result,err = json.Marshal(list.contentListMsg)
	if err != nil {
		fmt.Println("marshal failed,err is",err)
		return
	}
	return
}



func SendMsg(a dataCollect,b... string){

	res := a.CollectMsg(b[1])
	fmt.Println("result is",string(res))
}

var a chan []interface{}
var wg sync.WaitGroup

func main() {

	test := titleListMsg{}

	test.CollectMsg("1234")

	fmt.Println("test is",test.titleListMsg)

	//test2 := contentListMsg{}

	//CollectMsg("1")b := test2.
	//
	//fmt.Println("test2 is",test2,b)

	fmt.Println("jiekoujieguo:")

	SendMsg(&test,"1234")

	
}
