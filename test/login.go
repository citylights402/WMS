package common

import (
	"database/sql"
	"errors"

	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "sam404"
	dbname   = "dbwms"
)

type User struct {
	uid      uint
	uname    string
	password string
}
type UserDao struct{} //初始化对象

//database/sql，实现对数据库 的连接（sql.Open）、查询（ db.Query）、操作（db.Exec）
var db *sql.DB // 连接池对象

func (ud *UserDao) OpenSql() (err error) {
	connInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// DB 代表一个具有零到多个底层连接的连接池，可以安全的被多个go程序同时使用
	//这里的open函数只是验证参数是否合法，而不会创建和数据库的连接,也不会检查账号密码是否正确
	db, err = sql.Open("postgres", connInfo)
	checkErr(err)

	err = db.Ping()
	checkErr(err)

	db.SetMaxOpenConns(20) //设置数据库连接池最大连接数
	db.SetMaxIdleConns(10) //设置最大空闲连接数

	fmt.Println("Successfully connected!")
	return nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//查询数据
func (ud *UserDao) QueryAll() (err error, userlist []User) {

	rows, err := db.Query("SELECT * FROM w_user")
	checkErr(err)

	userlist = make([]User, 0)
	var user User
	for rows.Next() {

		err = rows.Scan(&user.uid, &user.uname, &user.password)
		userlist = append(userlist, user)
		checkErr(err)
		fmt.Println("uid = ", user.uid, "\nuname = ", user.uname, "\npassword = ", user.password, "\n-----------")
	}
	if len(userlist) > 0 {
		for _, value := range userlist {
			fmt.Println(value.uid, value.uname, value.password)
		}
	} else {
		fmt.Println("数据库里没有数据")
		return errors.New("No such data exists in database"), userlist
	}
	return err, userlist
}

func (ud *UserDao) InsertMulti(userlist []User) (err error) {
	//插入数据
	stmt, err := db.Prepare("INSERT INTO w_user(uname,password) VALUES($1,$2) RETURNING uid")
	checkErr(err)

	var count int64
	for _, user := range userlist {
		res, err := stmt.Exec(user.uname, user.password)
		//这里的二个参数就是对应上面的$1,$2了

		checkErr(err)

		affect, err := res.RowsAffected()
		checkErr(err)
		count += affect
	}
	fmt.Println("rows affect:", count)
	return err
}
func (ud *UserDao) DeleteSql() {
	//删除数据
	stmt, err := db.Prepare("delete from w_user where uid=$1")
	checkErr(err)

	res, err := stmt.Exec(1)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("rows affect:", affect)
}

func (ud *UserDao) UpdateSql() {
	//更新数据
	stmt, err := db.Prepare("update w_user set uname=$1 where uid=$2")
	checkErr(err)
	// if  {

	// }
	res, err := stmt.Exec("ls", 2)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("rows affect:", affect)
}
func (ud *UserDao) CloseSql() {
	db.Close()
}

func sqlTest() {
	ud := new(UserDao)
	sep := "----------\n"
	ud.OpenSql()
	println(sep, "*sqlOpen")

	// ud.QueryAll()

	userlist := []User{{0, "zs", "222"}, {0, "ls", "333"}}
	ud.InsertMulti(userlist)

	// ud.DeleteSql()

	// ud.UpdateSql()

	defer ud.CloseSql()
	println(sep, "*sqlClose")
}

// func main() {
// 	sqlTest()
// }
