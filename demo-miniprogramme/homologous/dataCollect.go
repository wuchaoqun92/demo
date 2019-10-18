package main

import (
	"demo-person/demo-miniprogramme/common"
	"demo-person/demo-upload/upload"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
)

func DatabaseConnect() (db *sqlx.DB, err error) {
	database, err := sqlx.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/miniProgramme")
	if err != nil {
		fmt.Println("database open failed,err:", err)
		return
	}

	db = database
	return
}

type dataCollect interface {
	CollectMsg(string) []byte
}

type titleList struct {
	Id     int    `db:"list_id" json:"id"`
	Title  string `db:"list_title" json:"title"`
	Code   string `db:"code" json:"code"`
	Counts int    `db:"counts" json:"counts"`
}
type titleListMsg struct {
	titleListMsg []titleList
}

type contentList struct {
	ListId      int    `db:"list_id" json:"list_id"`
	Title       string `db:"title" json:"title"`
	Text        string `db:"text" json:"text"`
	ContentTime string `db:"content_time" json:"content_time"`
}
type contentListMsg struct {
	contentListMsg []contentList
}

type contentInsert struct {
	UserId int    `json:"user_id" db:"user_id"`
	ListId int    `json:"list_id" db:"list_id"`
	Titlt  string `json:"titlt" db:"titlt"`
	Text   string `json:"text" db:"text"`
}

func (list *titleListMsg) CollectMsg(code string) (result []byte) {

	db, err := DatabaseConnect()
	if err != nil {
		fmt.Println("database connect failed,err:", err)
		return
	}
	defer db.Close()

	str := `SELECT code,list_title,list_id,COUNT(content_id) as counts FROM
			(
			SELECT a.code,a.list_id,a.list_title,b.content_id
			FROM list a LEFT JOIN content b ON a.list_id = b.list_id
			WHERE a.code = ? 
					) c 
			GROUP BY list_id,code,list_title
			ORDER BY list_id`

	err = db.Select(&list.titleListMsg, str, code)
	if err != nil {
		fmt.Println("数据插入结构体失败，err", err)
		return
	}

	fmt.Println("gotdata is", list.titleListMsg)
	result, err = json.Marshal(list.titleListMsg)
	if err != nil {
		fmt.Println("marshal failed,err is", err)
		return
	}
	return
}

func (list *contentListMsg) CollectMsg(listId string) (result []byte) {

	db, err := DatabaseConnect()
	if err != nil {
		fmt.Println("database connect failed,err:", err)
		return
	}
	defer db.Close()

	err = db.Select(&list.contentListMsg, "SELECT list_id,text,title,content_time FROM content where list_id = ? ", listId)
	if err != nil {
		fmt.Println("数据插入结构体失败，err:", err)
		return
	}
	fmt.Println("gotdata is", list)
	result, err = json.Marshal(list.contentListMsg)
	if err != nil {
		fmt.Println("marshal failed,err is", err)
		return
	}
	return
}

func (content *contentInsert) CollectMsg(con string) (result []byte) {
	//将获取到的数据结构化
	err := json.Unmarshal([]byte(con), content)
	if err != nil {
		fmt.Println("unmarshal failed,err is", err)
		return
	}
	fmt.Println(content)

	//实例化返回数据协议
	backMsg := common.BackMessage{
		Cmd: "1",
		Msg: "success",
	}

	//链接数据库
	db, err := DatabaseConnect()
	if err != nil {
		fmt.Println("database connect failed,err:", err)
		return
	}
	defer db.Close()

	//检查数据库是否存在list_id,不存在则无法插入数据
	counts := 0
	row := db.QueryRow("SELECT count(*) FROM list WHERE list_id =? ", content.ListId)
	err = row.Scan(&counts)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("row:", counts)

	if counts == 1 { //list_id存在，进行插入操作
		//插入数据
		_, err = db.Exec("INSERT INTO content (list_id,user_id,title,text,create_time) VALUES (?,?,?,?,?)", content.ListId, content.UserId, content.Titlt, content.Text, time.Now().Format("2006-01-02 15:04:05"))
		if err != nil {
			fmt.Println("insert failed ,err is", err)
			backMsg.Cmd = "0"
			backMsg.Msg = fmt.Sprintf("%s", err)
			return
		}
	}

	if counts == 0 {
		err = errors.New("分组标题不存在")
		backMsg.Msg = fmt.Sprintf("%s", err)
	}

	result, err = json.Marshal(backMsg)
	return
}

func SendMsg(a dataCollect, b string) {

	res := a.CollectMsg(b)
	fmt.Println("result is", string(res))
}

var a chan []interface{}
var wg sync.WaitGroup

func main() {

	test := titleListMsg{}

	test.CollectMsg("1234")

	fmt.Println("test is", test.titleListMsg)

	//test2 := contentListMsg{}

	//CollectMsg("1")b := test2.
	//
	//fmt.Println("test2 is",test2,b)

	fmt.Println("jiekoujieguo:")

	SendMsg(&test, "1234")

	x := contentInsert{
		UserId: 1,
		ListId: 111,
		Titlt:  "lalala",
		Text:   "阿杀得快解放；拉时间对方； 啊；了都是费劲；卡点十几分；爱空间撒旦法；卡萨丁；分卡大富科技阿是看京东方",
	}

	y, _ := json.Marshal(x)

	SendMsg(&x, string(y))

	//upload.Upload("http://localhost:8888/upload","/Users/wuchaoqun/Documents/Work/缺陷跟踪/东莞第一鸡_1.mp4")
	upload.Upload("http://localhost:8888/upload", "/Users/wuchaoqun/Documents/Work/缺陷跟踪/1m.mp4")

}
