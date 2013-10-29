/*
 *@author 菠菜君
 *@Version 0.2
 *@time 2013-10-29
 *@go语言实现微信公众平台
 *@青岛程序员 微信订阅号	qdprogrammer
 *@Golang 微信订阅号	gostock
 *@关于青岛程序员的技术，创业，生活 分享。
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

const (
	TOKEN    = "gostock"
	Text     = "text"
	Location = "location"
	Image    = "image"
	Link     = "link"
	Event    = "event"
	Music    = "music"
	News     = "news"
)

type msgBase struct {
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
}

type Request struct {
	XMLName                xml.Name `xml:"xml"`
	msgBase                         // base struct
	Location_X, Location_Y float32
	Scale                  int
	Label                  string
	PicUrl                 string
	MsgId                  int
}

type Response struct {
	XMLName xml.Name `xml:"xml"`
	msgBase
	ArticleCount int     `xml:",omitempty"`
	Articles     []*item `xml:"Articles>item,omitempty"`
	FuncFlag     int
}

type item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string
	Description string
	PicUrl      string
	Url         string
}

func weixinEvent(w http.ResponseWriter, r *http.Request) {
	if weixinCheckSignature(w, r) == false {
		fmt.Println("auth failed, attached?")
		return
	}

	fmt.Println("auth success, parse POST")

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(string(body))
	var wreq *Request
	if wreq, err = DecodeRequest(body); err != nil {
		log.Fatal(err)
		return
	}

	wresp, err := dealwith(wreq)
	if err != nil {
		log.Fatal(err)
		return
	}

	data, err := wresp.Encode()
	if err != nil {
		fmt.Printf("error:%v\n", err)
		return
	}

	fmt.Println(string(data))
	fmt.Fprintf(w, string(data))
	return
}

func dealwith(req *Request) (resp *Response, err error) {
	resp = NewResponse()
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	resp.MsgType = Text

	if req.MsgType == Event {
		if req.Content == "subscribe" {
			resp.Content = "欢迎关注微信订阅号qdprogrammer, 分享青岛程序员的技术，创业，生活。"
			return resp, nil
		}
	}

	if req.MsgType == Text {
		if strings.Trim(strings.ToLower(req.Content), " ") == "help" {
			resp.Content = "欢迎关注微信订阅号qdprogrammer, 分享青岛程序员的技术，创业，生活。"
			return resp, nil
		}
		resp.Content = "亲，菠菜君已经收到您的消息, 将尽快回复您."
	} else {
		resp.Content = "暂时还不支持其他的类型"
	}
	return resp, nil
}

func weixinAuth(w http.ResponseWriter, r *http.Request) {
	if weixinCheckSignature(w, r) == true {
		var echostr string = strings.Join(r.Form["echostr"], "")
		fmt.Fprintf(w, echostr)
	}
}

func weixinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("GET begin...")
		weixinAuth(w, r)
		fmt.Println("GET END...")
	} else {
		fmt.Println("POST begin...")
		weixinEvent(w, r)
		fmt.Println("POST END...")
	}
}

func main() {
	http.HandleFunc("/check", weixinHandler)
	//http.HandleFunc("/", action)
	port := "80"
	println("Listening on port ", port, "...")
	err := http.ListenAndServe(":"+port, nil) //设置监听的端口

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func str2sha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func weixinCheckSignature(w http.ResponseWriter, r *http.Request) bool {
	r.ParseForm()
	fmt.Println(r.Form)

	var signature string = strings.Join(r.Form["signature"], "")
	var timestamp string = strings.Join(r.Form["timestamp"], "")
	var nonce string = strings.Join(r.Form["nonce"], "")
	tmps := []string{TOKEN, timestamp, nonce}
	sort.Strings(tmps)
	tmpStr := tmps[0] + tmps[1] + tmps[2]
	tmp := str2sha1(tmpStr)
	if tmp == signature {
		return true
	}
	return false
}

func DecodeRequest(data []byte) (req *Request, err error) {
	req = &Request{}
	if err = xml.Unmarshal(data, req); err != nil {
		return
	}
	req.CreateTime *= time.Second
	return
}

func NewResponse() (resp *Response) {
	resp = &Response{}
	resp.CreateTime = time.Duration(time.Now().Unix())
	return
}

func (resp Response) Encode() (data []byte, err error) {
	resp.CreateTime = time.Second
	data, err = xml.Marshal(resp)
	return
}
