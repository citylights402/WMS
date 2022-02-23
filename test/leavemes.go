package leavemes

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type Message struct {
	Id         int
	Uname      string
	Content    string
	InsertTime string
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "sam404"
	dbname   = "dbwms"
)

func CheckErr(err error) {
	if err != nil {
		fmt.Println("11111")
		panic(err)
	}
}
func MessageList(w http.ResponseWriter, r *http.Request) {
	var id int
	var insert_time int
	var uname string
	var content string

	connInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connInfo)
	CheckErr(err)
	rows, err := db.Query("select uid, uname,message_content,  message_time from w_message")
	CheckErr(err)
	var msgSlice []*Message
	for rows.Next() {
		err = rows.Scan(&id, &uname, &content, &insert_time)
		CheckErr(err)
		msg := new(Message)
		msg.Id = id
		msg.InsertTime = time.Unix(int64(insert_time), 0).Format("2006-01-02 15:04:05")
		msg.Uname = uname
		msg.Content = content

		//fmt.Fprintf(w, id)
		msgSlice = append(msgSlice, msg)
	}

	//解析到模板
	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, &msgSlice)

}
func Add(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	uname := "hahahah 没有"
	content := r.FormValue("content")
	insertTime := time.Now().Unix()

	db, err := sql.Open("mysql", "root:123456@tcp(172.20.10.80:3306)/test?charset=utf8")
	CheckErr(err)
	stmt, err := db.Prepare("insert into book(uname,content,insert_time)values(?,?,?)")
	if _, err := stmt.Exec(uname, content, insertTime); err == nil {
		//w.Write([]byte("ok"))
	}

	MessageList(w, r)
	defer db.Close()
}

// func main() {
// 	http.HandleFunc("/", MessageList)
// 	http.HandleFunc("/list", MessageList)
// 	http.HandleFunc("/add", add)
// 	http.Handle("/static/", http.FileServer(http.Dir("./")))
// 	//http.HandleFunc("/", sayhelloName)       //设置访问的路由
// 	err := http.ListenAndServe(":9090", nil) //设置监听的端口
// 	if err != nil {
// 		log.Fatal("ListenAndServe: ", err)
// 	}
// }
