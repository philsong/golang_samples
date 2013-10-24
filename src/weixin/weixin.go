/*
 *@author widuu
 *@time 2013-7-19
 *@go语言实现公众平台
 */
package main

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Request struct {
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
	MsgId        int
}

type Response struct {
	ToUserName   string
	FromUserName string `xml:"xml>FromUserName"`
	CreateTime   time.Duration
	MsgType      string `xml:"xml>MsgType"`
	Content      string `xml:"xml>Content"`
}

func str2sha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func action(w http.ResponseWriter, r *http.Request) {
	postedMsg, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	r.Body.Close()
	v := Request{}
	xml.Unmarshal(postedMsg, &v)
	fmt.Println(v)
	if v.MsgType == "text" {
		v := Response{v.FromUserName, v.ToUserName, time.Second, v.MsgType, v.Content}
		output, err := xml.Marshal(v)
		if err != nil {
			fmt.Printf("error:%v\n", err)
		}
		fmt.Println(string(output))
		fmt.Fprintf(w, string(output))
	} else if v.MsgType == "event" {
		Content := `"欢迎关注
								我的微信"`
		v := Response{v.ToUserName, v.FromUserName, time.Second, v.MsgType, Content}
		output, err := xml.Marshal(v)
		if err != nil {
			fmt.Printf("error:%v\n", err)
		}
		fmt.Println(string(output))
		fmt.Fprintf(w, string(output))
	}
}

func checkSignature(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()
		var token string = "gostock"
		var signature string = strings.Join(r.Form["signature"], "")
		var timestamp string = strings.Join(r.Form["timestamp"], "")
		var nonce string = strings.Join(r.Form["nonce"], "")
		var echostr string = strings.Join(r.Form["echostr"], "")
		tmps := []string{token, timestamp, nonce}
		sort.Strings(tmps)
		tmpStr := tmps[0] + tmps[1] + tmps[2]
		tmp := str2sha1(tmpStr)
		if tmp == signature {
			fmt.Fprintf(w, echostr)
		}
	} else {
		action(w, r)
	}

}

func main() {
	http.HandleFunc("/check", checkSignature)
	//http.HandleFunc("/", action)
	port := "80"
	println("Listening on port ", port, "...")
	err := http.ListenAndServe(":"+port, nil) //设置监听的端口

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
