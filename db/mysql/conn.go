package mysql

import (
	"database/sql"
	"fmt"
	_"github.com/go-sql-driver/mysql"
	"os"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:123456@tcp(192.168.59.128:13306)/fileserver?charset=utf8&parseTime=true")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Println("Failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}

}

//返回数据库链接对象
func DBConn() *sql.DB {
	return db
}


