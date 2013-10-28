package main

import (
	"database/sql" //这包一定要引用
	"encoding/json"
	"fmt"                              //这个前面一章讲过
	_ "github.com/go-sql-driver/mysql" //这就是刚才下载的包
)

// 定义一个结构体, 需要大写开头哦, 字段名也需要大写开头哦, 否则json模块会识别不了
// 结构体成员仅大写开头外界才能访问
type User struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
}

// 一如既往的main方法
func main() {
	// 格式有点怪, @tcp 是指网络协议(难道支持udp?), 然后是域名和端口
	db, e := sql.Open("mysql", "root:@tcp(192.168.7.15:3306)/5000w?charset=utf8")
	if e != nil { //如果连接出错,e将不是nil的
		print("ERROR?")
		return
	}

	defer db.Close()

	// 提醒一句, 运行到这里, 并不代表数据库连接是完全OK的, 因为发送第一条SQL才会校验密码 汗~!
	rows, e := db.Query("SELECT Address,CtfId,Gender FROM `最后5000` limit 1")
	if e != nil {
		fmt.Printf("query error!!%v\n", e)
		return
	}
	if rows == nil {
		print("Rows is nil")
		return
	}
	fmt.Println("DB rows.Next")
	for rows.Next() { //跟java的ResultSet一样,需要先next读取
		user := new(User)
		// rows貌似只支持Scan方法 继续汗~! 当然,可以通过GetColumns()来得到字段顺序
		row_err := rows.Scan(&user.User, &user.Password, &user.Host)
		if row_err != nil {
			print("Row error!!")
			return
		}
		b, _ := json.Marshal(user)
		fmt.Println(string(b)) // 这里没有判断错误, 呵呵, 一般都不会有错吧
	}
	fmt.Println("Done")
}
