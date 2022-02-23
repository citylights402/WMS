package common

import (
	"database/sql"
	"errors"
	"time"

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

type Message struct {
	Id         int
	Content    string
	InsertTime time.Time
	UId        int
}
type MessageDao struct{} //初始化对象

//database/sql，实现对数据库 的连接（sql.Open）、查询（ db.Query）、操作（db.Exec）
var db *sql.DB // 连接池对象

func (md *MessageDao) OpenSql() (err error) {
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

func (md *MessageDao) CloseSql() {
	db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//查询数据
func (md *MessageDao) QueryAll() (err error, messagelist []Message) {
	seq := "----------\n"
	md.OpenSql()
	println(seq, "*sqlOpen")

	rows, err := db.Query("SELECT * FROM w_message")
	checkErr(err)

	messagelist = make([]Message, 0)
	var message Message
	for rows.Next() {

		err = rows.Scan(&message.Id, &message.Content, &message.InsertTime, &message.UId)
		messagelist = append(messagelist, message)
		checkErr(err)
		fmt.Println("uid = ", message.Id, "\nuname = ", message.Content, "\npassword = ", message.InsertTime, "\n-----------")
	}
	if len(messagelist) > 0 {
		for _, message := range messagelist {
			fmt.Println(message.Id, message.Content, message.InsertTime)
		}
	} else {
		fmt.Println("数据库里没有数据")
		return errors.New("No such data exists in database"), messagelist
	}

	defer md.CloseSql()
	println(seq, "*sqlClose")

	return err, messagelist
}

func (md *MessageDao) InsertMulti(messagelist []Message) (err error) {
	seq := "----------\n"
	md.OpenSql()
	println(seq, "*sqlOpen")

	//插入数据
	stmt, err := db.Prepare("INSERT INTO w_message(message_content,message_time,user_id) VALUES($1,$2,$3) RETURNING id")
	checkErr(err)

	var count int64

	for _, message := range messagelist {
		// temp := time.Unix(int64(message.InsertTime), 0).Format("2006-01-02 15:04:05")
		res, err := stmt.Exec(message.Content, message.InsertTime, message.UId)
		//这里的二个参数就是对应上面的$1,$2了

		checkErr(err)

		affect, err := res.RowsAffected()
		checkErr(err)
		count += affect
	}
	fmt.Println("rows affect:", count)

	defer md.CloseSql()
	println(seq, "*sqlClose")
	return err
}

func (md *MessageDao) InsertOne(message *Message) (err error) {
	seq := "----------\n"
	md.OpenSql()
	println(seq, "*sqlOpen")

	//插入数据
	stmt, err := db.Prepare("INSERT INTO w_message(message_content,message_time,user_id) VALUES($1,$2,$3) RETURNING id")
	checkErr(err)

	// temp := time.Unix(int64(message.InsertTime), 0).Format("2006-01-02 15:04:05")
	_, err = stmt.Exec(message.Content, message.InsertTime, message.UId)
	//这里的二个参数就是对应上面的$1,$2了

	checkErr(err)

	defer md.CloseSql()
	println(seq, "*sqlClose")
	return err
}

func (md *MessageDao) DeleteOne(mid int) (err error) {
	seq := "----------\n"
	md.OpenSql()
	println(seq, "*sqlOpen")

	//删除数据
	stmt, err := db.Prepare("delete from w_message where id=$1")
	checkErr(err)

	res, err := stmt.Exec(mid)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("rows affect:", affect)
	defer md.CloseSql()
	println(seq, "*sqlClose")
	return err
}

// func (ud *UserDao) UpdateSql() {
// 	//更新数据
// 	stmt, err := db.Prepare("update w_user set uname=$1 where uid=$2")
// 	checkErr(err)
// 	// if  {

// 	// }
// 	res, err := stmt.Exec("ls", 2)
// 	checkErr(err)

// 	affect, err := res.RowsAffected()
// 	checkErr(err)

// 	fmt.Println("rows affect:", affect)
// }

func SqlTest() {
	md := new(MessageDao)
	sep := "----------\n"
	md.OpenSql()
	println(sep, "*sqlOpen")

	// md.QueryAll()
	// temp := time.Unix(int64(message.InsertTime), 0).Format("2006-01-02 15:04:05")

	messagelist := []Message{{0, "留言", time.Now(), 2}, {0, "留言3", time.Now(), 3}}
	md.InsertMulti(messagelist)

	// ud.DeleteSql()

	// ud.UpdateSql()

	defer md.CloseSql()
	println(sep, "*sqlClose")
}

// func main() {
// 	sqlTest()
// }
