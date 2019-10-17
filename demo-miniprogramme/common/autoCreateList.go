package common

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jmoiron/sqlx"
)

func CreateList(title string) {
	db, err := DatabaseConnect()
	if err != nil {
		fmt.Println("database connect failed,err:", err)
		return
	}
	defer db.Close()

	str := `SELECT DISTINCT code FROM list`

	code := CommonTodo(str)
	fmt.Println(code)
	for _, v := range code {
		db.Exec("INSERT INTO list (list_title,code) VALUES (?,?)", title, v)
	}

}

func DeleteList() {
	db, err := DatabaseConnect()
	if err != nil {
		fmt.Println("database connect failed,err:", err)
		return
	}
	defer db.Close()

	str := `SELECT list_id FROM
			(
			SELECT a.code,a.list_id,a.list_title,b.content_id
			FROM list a LEFT JOIN content b ON a.list_id = b.list_id
			WHERE a.code = ?
					) c 
			GROUP BY list_id having COUNT(content_id) = 0
			ORDER BY list_id`

	code := CommonTodo(str)
	for _, v := range code {
		fmt.Println(v)
		db.Exec("DELETE from list where list_id = ?", v)
	}
}

func CommonTodo(str string) (code []string) {
	db, err := DatabaseConnect()
	if err != nil {
		fmt.Println("database connect failed,err:", err)
		return
	}
	defer db.Close()

	rows, err := db.Queryx(str)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	code = make([]string, 0)

	for rows.Next() {
		var a string
		err = rows.Scan(&a)
		code = append(code, a)
	}
	fmt.Println(code)
	return
}
