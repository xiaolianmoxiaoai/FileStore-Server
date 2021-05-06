package db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
)

func UserSignUp(username, password string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user(`user_name`, `user_pwd`)values(?,?)")
	if err != nil {
		fmt.Println("Failed to insert! err:" + err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return false
	}
	if rowsAffect, err := ret.RowsAffected(); nil == err && rowsAffect > 0 {
		fmt.Println("User sign up success!")
		return true
	} else {
		fmt.Println("User has already signed up!")
		return false
	}
}

func UserLoginIn(username, encpasswd string) bool {
	stmt, err := mydb.DBConn().Prepare("select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found:" + username)
		return false
	}
	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpasswd {
		return true
	}
	return false

}

//更新token到db
func UpdateToken(username, token string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"replace into tbl_user_token (`user_name`, `user_token`)values(?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

//获取用户token
func GetUserToken(username string) string {
	stmt, err := mydb.DBConn().Prepare("select user_token from tbl_user_token where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		return ""
	} else if rows == nil {
		return ""
	}
	pRows := mydb.ParseRows(rows)
	return string(pRows[0]["user_token"].([]byte))
}

type User struct {
	Username string
	Email string
	Phone string
	SignupAt string
	LastActiveAt string
	status int
}
//获取用户信息
func GetUserInfo(username string) (User, error) {
	user := User{}
	stmt, err := mydb.DBConn().Prepare("select user_name, signup_at from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()
	stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil

}