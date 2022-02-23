package main

import (
	"database/sql"

	"fmt"

	_ "github.com/lib/pq"
)

//database/sql，实现对数据库 的连接（sql.Open）、查询（ db.Query）、操作（db.Exec）
var db *sql.DB

func sqlOpen() {
	var err error
	db, err = sql.Open("postgres", "port=5432 user=postgres password=sam404 dbname=dbwms sslmode=disable") //sslmode就是安全验证模式
	checkErr(err)
	err = db.Ping()
	checkErr(err)

	fmt.Println("Successfully connected!")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func sqlInsert() {
	//插入数据
	stmt, err := db.Prepare("INSERT INTO w_user(uname,password) VALUES($1,$2) RETURNING uid")
	checkErr(err)

	res, err := stmt.Exec("zs", "222")
	//这里的二个参数就是对应上面的$1,$2了

	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("rows affect:", affect)
}
func sqlDelete() {
	//删除数据
	stmt, err := db.Prepare("delete from w_user where uid=$1")
	checkErr(err)

	res, err := stmt.Exec(1)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("rows affect:", affect)
}
func sqlSelect() {
	//查询数据
	rows, err := db.Query("SELECT * FROM w_user")
	checkErr(err)

	println("-----------")
	for rows.Next() {
		var uid int
		var uname string
		var password string
		err = rows.Scan(&uid, &uname, &password)
		checkErr(err)
		fmt.Println("uid = ", uid, "\nuname = ", uname, "\npassword = ", password, "\n-----------")
	}
}
func sqlUpdate() {
	//更新数据
	stmt, err := db.Prepare("update w_user set uname=$1 where uid=$2")
	checkErr(err)

	res, err := stmt.Exec("ls", 2)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("rows affect:", affect)
}
func sqlClose() {
	db.Close()
}

func sqlTest() {

	sep := "----------\n"
	sqlOpen()
	println(sep, "*sqlOpen")

	// sqlSelect()
	// println(sep, "*sqlSelect")

	// sqlInsert()
	// sqlSelect()
	// println(sep, "*sqlInsert")

	sqlUpdate()
	sqlSelect()
	println(sep, "*sqlUpdate")

	// sqlDelete()
	// sqlSelect()
	// println(sep, "*sqlDelete")

	defer sqlClose()
	println(sep, "*sqlClose")
}

func main() {
	// common.OpenSql()
	// sqlTest()
}
