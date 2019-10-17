package common

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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
