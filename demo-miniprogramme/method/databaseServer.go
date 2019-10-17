package method

import (
	"demo-person/demo-miniprogramme/common"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

//type content_List struct {
//	Code string `db:"code"`
//	CotentId int `db:"contentid"`
//	ListTitle string `db:"listTitle"`
//}
//
//type content struct {
//	content_List
//	Title string `db:"title"`
//	Author string `db:"userid"`
//	CreateTime string `db:"createTime"`
//	Text string `db:"text"`
//}

func CheckCode(code, openid string) (err error) {
	db, err = common.DatabaseConnect()
	if err != nil {
		fmt.Println("database connect failed,err:", err)
		return
	}
	defer db.Close()

	var counts int

	//判断当前用户是否两者都正确
	row := db.QueryRow("SELECT count(*) FROM users WHERE code =? and openid=?", code, openid)
	err = row.Scan(&counts)
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println("row:", counts)
	if counts > 0 {
		return
	}

	//满足两者的信息不存在时，检查验证码是否正确，创建新用户
	row = db.QueryRow("SELECT count(*) FROM users WHERE code =?", code)
	err = row.Scan(&counts)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	fmt.Println("row:", counts)
	if counts > 0 {
		//验证码正确时，检查用户是否已存在
		row := db.QueryRow("SELECT count(*) FROM users WHERE openid =?", openid)
		err = row.Scan(&counts)
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		fmt.Println(counts)

		if counts == 0 {
			_, err = db.Exec("INSERT INTO users (openid) VALUES (?)", openid) //插入数据
			if err != nil {
				fmt.Println("err:", err)
				return
			}
		} else {
			err = errors.New("用户已存在，请使用正确的验证码")
			fmt.Println(err)
			return
		}
	} else {
		err = errors.New("验证码不存在")
		fmt.Println(err)
		return
	}
	return
}

func GetListFromDatabase(code string) (list []Content) {
	db, err := common.DatabaseConnect()
	if err != nil {
		fmt.Println("database connect failed,err:", err)
		return
	}
	defer db.Close()

	list = make([]Content, 0)
	str := `SELECT code,list_title,list_id,COUNT(content_id) as counts FROM
			(
			SELECT a.code,a.list_id,a.list_title,b.content_id
			FROM list a LEFT JOIN content b ON a.list_id = b.list_id
			WHERE a.code = ? 
					) c 
			GROUP BY list_id,code,list_title
			ORDER BY list_id`

	err = db.Select(&list, str, code)
	//err = db.Select(&list, "SELECT code,list_id,list_title FROM list  where code = ? ORDER BY list_id DESC",code)
	if err != nil {
		fmt.Println("数据插入结构体失败，err", err)
		return
	}
	fmt.Println("gotdata is", list)
	return
}

func GetContentDetailFromDatabase(content_id string) (list []Content) {
	db, err := common.DatabaseConnect()
	if err != nil {
		fmt.Println("database connect failed,err:", err)
		return
	}
	defer db.Close()

	list = make([]Content, 0)

	err = db.Select(&list, "SELECT title,user_id,create_time,text,content_time FROM content  where content_id = ? ", content_id)
	if err != nil {
		fmt.Println("数据插入结构体失败，err", err)
		return
	}
	fmt.Println("gotdata is", list)
	return
}

func InsertContent() (err error) {

	_, err = db.Exec("INSERT INTO content (title,userid,createTime,text) VALUES (?,?,?)", 100, "wuchaoqun", 18)

	return
}

//err := db.Select(&stu, "SELECT * FROM stu  where id = ? ORDER BY id ASC",1)  //将查询结果插入stu结构体

//db.Exec("INSERT INTO stu (id, name, age) VALUES (?,?,?)",100,"wuchaoqun",18)  //插入数据

//db.Exec("UPDATE stu set name=? where id = ?",100,"henlihai")  //更新数据

//db.Exec("DELETE from stu where id = ?",100)  //删除数据

/* 获取列表，及对应列下内容条数
SELECT code,list_title,list_id,COUNT(content_id) FROM
(
SELECT a.code,a.list_id,a.list_title,b.content_id
FROM list a LEFT JOIN content b ON a.list_id = b.list_id
WHERE a.code = '1234'
		) c
GROUP BY code,list_title,list_id
*/
