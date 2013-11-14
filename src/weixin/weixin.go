/*
 *@author 菠菜
 *@Version 0.4
 *@Update time 2013-11-14
 *@golang 微信公众平台API引擎开发
 *@演示 微信订阅号	gostock
 *@开源 https://github.com/philsong/
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
	"stock"
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

type item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string
	Description string
	PicUrl      string
	Url         string
}

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

func weixinProc(w http.ResponseWriter, r *http.Request) {
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

func dealwithEvent(req *Request, resp *Response) {
	if req.Content == "subscribe" {
		resp.Content = "欢迎关注微信订阅号qdprogrammer, 分享青岛程序员的技术，创业，生活。"
	}
}

func dealwithText(keyword string, resp *Response) {

	if keyword == "help" {
		resp.Content = "欢迎关注微信订阅号qdprogrammer, 分享青岛程序员的技术，创业，生活。"
	} else if keyword == "会员卡" {
		//todo: bind mobile to mysql database.
		//if req.FromUserName in db, then alread bind
		//else req.FromUserName not in bd, then remind to bind
		resp.Content = "发送手机号绑定会员卡。"
	} else if keyword == "股票" {
		resp.Content = stock.Deal(keyword)
		//todo
	} else {
		resp.Content = "亲，菠菜君已经收到您的消息, 将尽快回复您."
	}
}

func dealwithImage(req *Request, resp *Response) {
	var a item
	a.Description = "雅蠛蝶。。。^_^^_^1024你懂的"
	a.Title = "雅蠛蝶图文测试"
	a.PicUrl = "http://static.yaliam.com/gwz.jpg"
	a.Url = "http://blog.csdn.net/songbohr"

	resp.MsgType = News
	resp.ArticleCount = 1
	resp.Articles = append(resp.Articles, &a)
	resp.FuncFlag = 1
}

func dealwith(req *Request) (resp *Response, err error) {
	resp = NewResponse()
	resp.ToUserName = req.FromUserName
	resp.FromUserName = req.ToUserName
	resp.MsgType = Text

	if req.MsgType == Event {
		dealwithEvent(req, resp)
	} else if req.MsgType == Text {
		keyword := strings.Trim(strings.ToLower(req.Content), " ")
		dealwithText(keyword, resp)
	} else if req.MsgType == Image {
		dealwithImage(req, resp)
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
		weixinAuth(w, r)
	} else {
		weixinProc(w, r)
	}
}

func main() {
	http.HandleFunc("/check", weixinHandler)
	port := "80"
	println("Weixin Listening on port ", port, "...")
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
	resp.CreateTime = time.Duration(time.Now().Unix())
	data, err = xml.Marshal(resp)
	return
}
