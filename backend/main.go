package main

import (
	"WMS/backend/common"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//留言
func ShowMsg(c *gin.Context) {
	// 从后台取取数据

	// c.JSON：方法一返回JSON格式的数据
	// c.JSON(200, gin.H{
	// 	"message": "POST",
	// })

	// 方法二：使用结构体
	md := new(common.MessageDao)
	err, messagelist := md.QueryAll()
	if err != nil {
		panic(err)
	}
	for _, message := range messagelist {
		c.JSON(http.StatusOK, message) //传给前端,并附上code标识符
	}

}
func PostMsg(c *gin.Context) {
	// 注意：下面为了举例子方便，暂时忽略了错误处理
	b, _ := c.GetRawData() // 从c.Request.Body读取请求数据
	// 定义map或结构体
	// var m map[string]interface{}
	md := new(common.MessageDao)
	m := new(common.Message)

	// 反序列化
	_ = json.Unmarshal(b, &m)
	fmt.Println(m)

	md.InsertOne(m)
	c.JSON(http.StatusOK, m)

}
func DeleteMsg(c *gin.Context) {
	res, _ := c.GetRawData()
	fmt.Println("------------------", res)
	md := new(common.MessageDao)
	id := make(map[string]int)
	_ = json.Unmarshal(res, &id)
	fmt.Println(id["MId"])
	md.DeleteOne(id["MId"])
	c.JSON(http.StatusOK, id)
}
func main() {
	// 创建一个默认的路由引擎
	r := gin.Default()

	// 查询
	r.GET("/showmsg", ShowMsg)

	//接受前台数据
	r.POST("/postmsg", PostMsg)

	//删除留言
	r.GET("/deletemsg", DeleteMsg)

	// 启动HTTP服务，默认在0.0.0.0:8080启动服务
	r.Run()

	//测试
	// common.SqlTest()

}
