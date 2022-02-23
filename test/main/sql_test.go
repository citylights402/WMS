package test

import (
	"database/sql" //通用的接口
	"fmt"
	"testing"

	_ "github.com/bmizerany/pq" //必须要有相应的驱动
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "sam404" //你自己数据库的密码
	dbname   = "wms"    //创建的数据库
)

var db *sql.DB // 连接池对象
var err error

type product struct {
	ProductNo string
	Name      string
	Price     float64
}

type productDao struct{}

func (pd *productDao) initDB() (err error) {
	pdqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	//  pdqlInfo := "postgres:密码@tcp(127.0.0.1:5432)/test" // 用户名:密码@tcp(ip端口)/数据库名字，暂时出错
	db, err = sql.Open("postgres", pdqlInfo) //Open(driverName 驱动名字, dataSourceName string 数据库信息)
	// DB 代表一个具有零到多个底层连接的连接池，可以安全的被多个go程序同时使用
	//这里的open函数只是验证参数是否合法，而不会创建和数据库的连接,也不会检查账号密码是否正确
	if err != nil {
		fmt.Println("Wrong args.Connected failed.")
		return err
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Connected failed.")
		return err
	}
	db.SetMaxOpenConns(20) //设置数据库连接池最大连接数
	db.SetMaxIdleConns(10) //设置最大空闲连接数
	fmt.Println("Successfully connected!")
	return nil
}

func (pd *productDao) closeDB() (err error) {
	return db.Close()
}

func (pd *productDao) doQueryAll() (error, []product) {
	rows, err := db.Query(`Select * from products`)
	if err != nil {
		fmt.Println("Some amazing wrong happens in the process of Query.", err)
		return err, []product{}
	}
	products := make([]product, 0)
	defer rows.Close() //关闭连接
	index := 0
	var p product
	for rows.Next() {
		err := rows.Scan(&p.ProductNo, &p.Name, &p.Price)
		products = append(products, p)
		if err != nil { // 获得的都是字符串
			fmt.Println("Some amazing wrong happens in the process of queryAll.", err)
			return err, products
		}
		index++
	}
	if index > 0 {
		fmt.Println("The data of table is as follow.")
		for _, p := range products {
			fmt.Printf("%v %s %v\n", p.ProductNo, p.Name, p.Price)
		}
		fmt.Println("Successfully query ", len(products))
		return nil, products
	} else {
		fmt.Println("No such data exists in database. ")
		return fmt.Errorf("No such data exists in database. "), products
	}
}

func (pd *productDao) doQueryByPrice(price float32) (error, []product) {
	rows, err := db.Query(`Select * from products where price=$1`, price)
	if err != nil {
		fmt.Println("Some amazing wrong happens in the process of Query.", err)
		return err, []product{}
	}
	products := make([]product, 0)
	defer rows.Close() //关闭连接
	index := 0
	var p product
	for rows.Next() {
		err := rows.Scan(&p.ProductNo, &p.Name, &p.Price)
		products = append(products, p)
		if err != nil { // 获得的都是字符串
			fmt.Println("Some amazing wrong happens in the process of queryAll.", err)
			return err, products
		}
		index++
	}
	if index > 0 {
		fmt.Println("The data of table is as follow.")
		for _, p := range products {
			fmt.Printf("%v %s %v\n", p.ProductNo, p.Name, p.Price)
		}
		fmt.Println("Successfully query ", len(products))
		return nil, products
	} else {
		fmt.Println("No such data exists in database. ")
		return fmt.Errorf("No such data exists in database. "), products
	}
}

func (pd *productDao) doQueryByName(pname string) (error, []product) {
	rows, err := db.Query(`Select * from products where name=$1`, pname)
	if err != nil {
		fmt.Println("Some amazing wrong happens in the process of Query.", err)
		return err, []product{}
	}
	products := make([]product, 0)
	defer rows.Close() //关闭连接
	index := 0
	var p product
	for rows.Next() {
		err := rows.Scan(&p.ProductNo, &p.Name, &p.Price)
		products = append(products, p)
		if err != nil { // 获得的都是字符串
			fmt.Println("Some amazing wrong happens in the process of queryAll.", err)
			return err, products
		}
		index++
	}
	if index > 0 {
		fmt.Println("The data of table is as follow.")
		for _, p := range products {
			fmt.Printf("%v %s %v\n", p.ProductNo, p.Name, p.Price)
		}
		fmt.Println("Successfully query ", len(products))
		return nil, products
	} else {
		fmt.Println("No such data exists in database. ")
		return fmt.Errorf("No data exists in database. "), products
	}
}

func (pd *productDao) doPreQueryByName(pname string) (error, []product) {
	sqlStr := "Select * from products where name=$1"
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare failed, err:%v\n", err)
		return err, []product{}
	}
	defer stmt.Close()
	rows, err := stmt.Query(pname)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return err, []product{}
	}
	defer rows.Close()
	index := 0
	products := make([]product, 0)
	var p product
	// 循环读取结果集中的数据
	for rows.Next() {
		err := rows.Scan(&p.ProductNo, &p.Name, &p.Price)
		products = append(products, p)
		if err != nil { // 获得的都是字符串
			fmt.Println("Some amazing wrong happens in the process of queryAll.", err)
			return err, []product{}
		}
		index++
	}
	if index > 0 {
		fmt.Println("The data of table is as follow.")
		for _, p := range products {
			fmt.Printf("%v %s %v\n", p.ProductNo, p.Name, p.Price)
		}
		fmt.Println("Successfully query ", len(products))
		return nil, products
	} else {
		fmt.Println("No such data exists in database. ")
		return fmt.Errorf("No data exists in database. "), products
	}
}

func (pd *productDao) doQueryByNo(pno string) (error, product) {
	var p product
	// 1、写查询单条记录的sql语句
	sqlStr := "Select * from  products where ProductNo = $1"
	// 2、执行
	row := db.QueryRow(sqlStr, pno)
	// 3、获得结果
	err = row.Scan(&p.ProductNo, &p.Name, &p.Price)
	if err != nil {
		fmt.Println("Some amazing wrong happens in the process of queryByNo :", err)
		return err, product{}
	}
	fmt.Printf("%v %s %v\n", p.ProductNo, p.Name, p.Price)
	fmt.Println("Successfully query 1")
	return nil, p
}

func (pd *productDao) doInsertMulti(projects []product) error {
	stmt, err := db.Prepare("Insert into products(Product_No,Name,Price) values($1,$2,$3)")
	if err != nil {
		fmt.Println("Some amazing wrong happens in preparation for the insert.")
		return err
	}
	defer stmt.Close()
	for i, project := range projects {
		_, err := stmt.Exec(project.ProductNo, project.Name, project.Price)
		if err != nil {
			fmt.Println("Some amazing wrong happens in the process of doInsertMulti.")
			fmt.Println("But successfully insert ", i)
			return err
		}
	}
	fmt.Println("Successfully add", len(projects))
	return nil
}

func (pd *productDao) doInsertOne(project product) error {
	stmt, err := db.Prepare("Insert into products(Product_No,Name,Price) values($1,$2,$3)")
	if err != nil {
		fmt.Println("Some amazing wrong happens in preparation for the insert.")
		//log.Fatal(err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(project.ProductNo, project.Name, project.Price)
	if err != nil {
		fmt.Println("Some amazing wrong happens in the process of doInsertOne.")
		//log.Fatal(err)
		return err
	}
	num, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Some amazing wrong happens in the affected fo deletion.")
		return err
	}
	fmt.Println("Successfully add ", num)
	return nil
}

func (pd *productDao) doDeleteByNo(pno string) error {
	stmt, err := db.Prepare(`Delete FROM products  where productno = $1`)
	if err != nil {
		fmt.Println("Some amazing wrong happens in preparation for the deletion.")
		return err
	}
	res, err := stmt.Exec(pno)
	if err != nil {
		fmt.Println("Some amazing wrong happens in execution for the deletion.")
		return err
	}
	num, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Some amazing wrong happens in the affected fo deletion.")
		return err
	}
	fmt.Println("Successfully delete", num)
	return nil
}

func (pd *productDao) doDeleteByName(pname string) error {
	stmt, err := db.Prepare(`Delete FROM products  where name = $1`)
	if err != nil {
		fmt.Println("Some amazing wrong happens in preparation for the deletion.")
		return err
	}
	res, err := stmt.Exec(pname)
	if err != nil {
		fmt.Println("Some amazing wrong happens in execution for the deletion.")
		return err
	}
	num, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Some amazing wrong happens in the affected fo deletion.")
		return err
	}
	fmt.Println("Successfully delete", num)
	return nil
}

func (pd *productDao) doUpdatePriceByNo(pno string, newValue float32) error {
	stmt, err := db.Prepare("Update products set price = $2 where  product_no = $1")
	if err != nil {
		fmt.Println("Some amazing wrong happens in preparation for the update.")
		//log.Fatal(err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(pno, newValue)
	if err != nil {
		fmt.Println("Some amazing wrong happens in execution for the update.")
		return err
	}
	num, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Some amazing wrong happens in the affected fo update.")
		return err
	}
	fmt.Println("Successfully update ", num)
	return nil
}

func (pd *productDao) doUpdatePriceByName(pname string, newValue float32) error {
	stmt, err := db.Prepare("Update products set price = $2 where  name = $1")
	if err != nil {
		fmt.Println("Some amazing wrong happens in preparation for the update.")
		//log.Fatal(err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(pname, newValue)
	if err != nil {
		fmt.Println("Some amazing wrong happens in execution for the update.")
		//log.Fatal(err)
		return err
	}
	num, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Some amazing wrong happens in the affected fo update.")
		return err
	}
	fmt.Println("Successfully update ", num)
	return nil
}

func TestSql(t *testing.T) {
	pd := new(productDao)
	err = pd.initDB()
	if err != nil {
		fmt.Println("initDB() failed. ")
	}
	defer pd.closeDB()
	pd.doQueryAll()
	pd.doPreQueryByName("apple")
}
